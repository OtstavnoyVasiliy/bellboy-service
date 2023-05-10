package db

import (
	"context"
	"database/sql"
	"fmt"
	"tg-bot/pkg/types"

	"github.com/jackc/pgx/v5"

	"github.com/spf13/viper"
)

type DataBase struct {
	Conn *pgx.Conn
}

type IDataBase interface {
	UserExists(ctx context.Context, user_id int) (bool, error)             // Check for user existing in database
	ActiveChat(ctx context.Context, group_id int64) (bool, error)          // Check for group existing in database
	GetAllGroups(ctx context.Context) ([]int64, error)                     // Get all bot's groups in database
	GetActiveGroups(ctx context.Context) ([]int64, error)                  // Get all bot's groups in database
	GetUser(ctx context.Context, user_id int) (string, string, error)      // Get user's first name and city by telegram username id
	CreateAuthSign(ctx context.Context, user_id int, tg_user_id int) error // Creating a relationship between b24 user id and telegram user id
	IsAuth(ctx context.Context, user_id int) (bool, error)                 // Check auth sign exist for current user
	RemoveGroup(ctx context.Context, group_id int64) error
	AddGroup(ctx context.Context, group_id int64) error
}

func NewDataBase(config *viper.Viper) (*DataBase, error) {
	connString := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", config.GetString("database.host"), config.GetInt("database.port"), config.GetString("database.name"), config.GetString("database.username"), config.GetString("database.password"))
	dbConf, err := pgx.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	conn, err := pgx.ConnectConfig(context.Background(), dbConf)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	return &DataBase{
		Conn: conn,
	}, nil
}

func (db *DataBase) UserExists(ctx context.Context, user_id int) (bool, error) {
	var exists bool
	err := db.Conn.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM tg_users WHERE id=$1)", user_id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (db *DataBase) ActiveChat(ctx context.Context, group_id int64) (bool, error) {
	var exists bool
	err := db.Conn.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM active_chats WHERE tg_chat=$1)", group_id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (db *DataBase) GetInviteLink(ctx context.Context, group_id int64) (string, error) {
	var link string
	if err := db.Conn.QueryRow(ctx, "SELECT invite_link FROM tg_chats WHERE id = $1;", group_id).Scan(&link); err != nil && err.Error() != "no rows in result set" {
		return "", err
	}
	return link, nil
}

func (db *DataBase) IsEmployee(ctx context.Context, user_id int) (bool, error) {
	var isEmployee bool
	err := db.Conn.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM bx_users JOIN tg_users ON bx_users.id = tg_users.bx_user WHERE tg_users.id = $1);", user_id).Scan(&isEmployee)
	if err != nil {
		return false, err
	}
	return isEmployee, nil
}

func (db  *DataBase) GetLeads(ctx context.Context) ([]int64, error) {
	var leads []int64

	rows, err := db.Conn.Query(ctx, "SELECT lead FROM bx_departments;")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var lead sql.NullInt64
		if err := rows.Scan(&lead); err != nil {
			return nil, err
		}
		if lead.Valid {
			leads = append(leads, lead.Int64)
		}
	}

	return leads, nil
}

func (db *DataBase) GetChatLinks(ctx context.Context, user_id int) (*types.ChatsResult, error) {
	var chats types.ChatsResult

	rows, err := db.Conn.Query(ctx, "SELECT all_ids, city_ids, dev_ids, rec_ids FROM get_chat_ids($1);", user_id)
	if err != nil && err.Error() != "no rows in result set" {
		return nil, err
	}

	for rows.Next() {
		var allID, cityID, devID, recID *types.ChatRecordShort
		err := rows.Scan(&allID, &cityID, &devID, &recID)
		if err != nil {
			return nil, err
		}

		if allID != nil {
			chats.All = append(chats.All, *allID)
		}
		if cityID != nil {
			chats.City = append(chats.City, *cityID)
		}
		if devID != nil {
			chats.Dev = append(chats.Dev, *devID)
		}
		if recID != nil {
			chats.Recommended = append(chats.Recommended, *recID)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &chats, nil
}

func (db *DataBase) AddGroup(ctx context.Context, group_id int64) error {
	if _, err := db.Conn.Exec(ctx, "INSERT INTO active_chats (tg_chat, joined_at) VALUES ($1, NOW());", group_id); err != nil {
		return err
	}

	return nil
}

func (db *DataBase) GetAllGroups(ctx context.Context) ([]int64, error) {
	var ids []int64
	rows, err := db.Conn.Query(ctx, "SELECT id FROM bot_groups;")
	if err != nil && err.Error() != "no rows in result set" {
		return nil, err
	}

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func (db *DataBase) GetActiveGroups(ctx context.Context) ([]int64, error) {
	var ids []int64
	rows, err := db.Conn.Query(ctx, "SELECT tg_chat FROM active_chats;")
	if err != nil && err.Error() != "no rows in result set" {
		return nil, err
	}

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func (db *DataBase) RemoveGroup(ctx context.Context, group_id int64) error {
	if _, err := db.Conn.Exec(ctx, "DELETE FROM active_chats WHERE tg_chat = $1;", group_id); err != nil {
		return err
	}
	return nil
}

func (db *DataBase) CreateAuthSign(ctx context.Context, user_id int, chat_id int) error {
	if _, err := db.Conn.Exec(ctx, "WITH updated_rows AS (UPDATE auth_signs SET auth_at = NOW() WHERE chat_id = $1 RETURNING * ) INSERT INTO auth_signs (chat_id, user_id, auth_at) SELECT $1, $2, NOW() WHERE NOT EXISTS (SELECT 1 FROM updated_rows);", chat_id, user_id); err != nil {
		return err
	}

	return nil
}

func (db *DataBase) IsAuth(ctx context.Context, user_id int) (bool, error) {
	var exists bool
	err := db.Conn.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM auth_signs WHERE chat_id=$1)", user_id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (db *DataBase) GetUser(ctx context.Context, user_id int) (string, string, error) {
	var firstName, city string
	err := db.Conn.QueryRow(ctx, "SELECT name, city FROM bx_users JOIN tg_users ON bx_users.id = tg_users.bx_user WHERE tg_users.id = $1;", user_id).Scan(&firstName, &city)
	if err != nil && err.Error() != "no rows in result set" {
		return "", "", err
	}
	return firstName, city, nil
}

func (db *DataBase) GetChatInfo(ctx context.Context, chat_id int64) (*types.ChatRecordShort, error) {
	var res types.ChatRecordShort

	if err := db.Conn.QueryRow(ctx, "SELECT title, description FROM tg_chats WHERE id=$1;", chat_id).Scan(&res.Title, &res.Description); err != nil && err.Error() != "no rows in result set" {
		return nil, err
	}

	return &res, nil
}

func (db *DataBase) GetBxInfo(ctx context.Context, user_id int) (int, string, string, error) {
	var leadBxID, leadTgID int
	var name, lastName string

	if err := db.Conn.QueryRow(ctx, "SELECT bx_departments.lead, bx_users.name, bx_users.last_name FROM bx_users INNER JOIN bx_departments ON bx_users.department = bx_departments.name WHERE bx_users.id = $1;", user_id).Scan(&leadBxID, &name, &lastName); err != nil && err.Error() != "no rows in result set" {
		return 0, "", "", err
	}

	if err := db.Conn.QueryRow(ctx, "SELECT id FROM tg_users WHERE bx_user = $1;", leadBxID).Scan(&leadTgID); err != nil && err.Error() != "no rows in result set" {
		return 0, "", "", err
	}

	return leadTgID, name, lastName, nil
}
