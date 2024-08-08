package handler

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	pb "github.com/xissg/file-system-grpc/internal"
	"github.com/xissg/file-system-grpc/utils"
	"log"
)

const (
	NormalUserStats = 0
)

type DBService struct {
	pb.UnimplementedDBServiceServer
	db *sqlx.DB
}

func NewDB() *DBService {
	dsn := "root:root@tcp(127.0.0.1:3306)/rpc_file_db"
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		panic(err)
		return nil
	}

	return &DBService{
		db: db,
	}
}
func (db *DBService) AddUser(_ context.Context, user *pb.User) (*pb.UserStatus, error) {
	id := utils.Snowflake()
	res := &pb.UserStatus{
		Success: false,
	}

	sqlStr := "INSERT INTO tbl_user(id,username,password,status) VALUES (?,?,?,?)"

	tx, err := db.db.Begin()
	if err != nil {
		log.Printf("failed to start tx error %v", err)
		return res, err
	}

	_, err = tx.Exec(sqlStr, id, user.Username, user.Password, NormalUserStats)

	if err != nil {
		log.Printf("failed to exec sql error %v", err)
		_ = tx.Rollback()
		return res, err
	}
	_ = tx.Commit()
	res.Success = true
	return res, nil
}
func (db *DBService) GetUser(_ context.Context, username *pb.UserName) (*pb.User, error) {

	user := User{}
	sqlStr := "SELECT id,username,password,status FROM tbl_user WHERE username = ?"
	err := db.db.Get(&user, sqlStr, username.Username)
	if err != nil {
		log.Printf("failed to get user error %v", err)
		return nil, err
	}

	res := &pb.User{
		Id:       user.ID,
		Username: user.UserName,
		Password: user.Password,
		Status:   user.Status,
	}
	return res, nil
}

func (db *DBService) CheckUserExist(_ context.Context, username *pb.UserName) (*pb.OK, error) {
	var exists bool
	sqlStr := "SELECT EXISTS(SELECT 1 FROM tbl_user WHERE username = ?)"
	err := db.db.Get(&exists, sqlStr, username.Username)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	if exists {
		return &pb.OK{
			Exist: true,
		}, nil
	} else {
		return &pb.OK{
			Exist: false,
		}, nil
	}
}
func (db *DBService) AddFile(_ context.Context, file *pb.File) (*pb.FileStatus, error) {
	id := utils.Snowflake()
	res := &pb.FileStatus{
		Success: false,
	}
	sqlStr := "INSERT INTO tbl_file (id,file_name,file_size, file_checksum,file_addr,user_id,status)VALUES (?,?,?,?,?,?,?)"
	tx, _ := db.db.Begin()
	_, err := tx.Exec(sqlStr, id, file.FileName, file.FileSize, file.Checksum, file.FileAddr, file.UserId, file.Status)
	if err != nil {
		log.Printf("failed to add file error %v", err)
		_ = tx.Rollback()
		return res, err
	}
	_ = tx.Commit()
	res.Success = true

	return res, nil
}

func (db *DBService) GetFileByUserID(_ context.Context, userID *pb.UserID) (*pb.Files, error) {
	var file []File
	sqlStr := "SELECT id,file_name,file_size,file_checksum,file_addr,user_id,status FROM tbl_file WHERE user_id = ?"
	err := db.db.Select(&file, sqlStr, userID.UserId)
	if err != nil {
		log.Printf("failed to get file by user id error %v", err)
		return nil, err
	}

	res := new(pb.Files)
	for _, f := range file {
		tmp := &pb.File{
			Id:       f.ID,
			FileName: f.FileName,
			FileSize: f.FileSize,
			Checksum: f.Checksum,
			FileAddr: f.FileAddr,
			Status:   f.Status,
			UserId:   f.UserId,
		}
		res.File = append(res.File, tmp)
	}

	return res, nil
}
func (db *DBService) GetFileByChecksum(_ context.Context, checksum *pb.Checksum) (*pb.File, error) {
	file := File{}
	sqlStr := "SELECT id,file_name,file_size,file_checksum,file_addr,user_id,status FROM tbl_file WHERE file_checksum = ?"
	err := db.db.Get(&file, sqlStr, checksum.Checksum)
	if err != nil {
		log.Printf("failed to get file by checksum error %v", err)
		return nil, err
	}

	res := &pb.File{
		Id:       file.ID,
		FileName: file.FileName,
		FileSize: file.FileSize,
		Checksum: file.Checksum,
		FileAddr: file.FileAddr,
		Status:   file.Status,
		UserId:   file.UserId,
	}
	return res, nil
}
func (db *DBService) UpdateFileStatus(_ context.Context, req *pb.UpdateRequest) (*pb.FileStatus, error) {
	res := &pb.FileStatus{
		Success: false,
	}
	sqlStr := "UPDATE tbl_file SET status = ? WHERE file_checksum = ?"
	tx, _ := db.db.Begin()
	_, err := tx.Exec(sqlStr, req.Status, req.Checksum)
	if err != nil {
		log.Printf("failed to update file status error %v", err)
		_ = tx.Rollback()
		return res, err
	}
	_ = tx.Commit()
	res.Success = true
	return res, nil
}
func (db *DBService) DeleteFile(_ context.Context, checksum *pb.Checksum) (*pb.FileStatus, error) {
	res := &pb.FileStatus{
		Success: false,
	}
	sqlStr := "DELETE FROM tbl_file WHERE file_checksum = ?"
	tx, _ := db.db.Begin()
	_, err := tx.Exec(sqlStr, checksum.Checksum)
	if err != nil {
		log.Printf("failed to delete file error %v", err)
		_ = tx.Rollback()
		return res, err
	}
	_ = tx.Commit()
	res.Success = true
	return res, nil
}

type User struct {
	ID       int64  `db:"id"`
	UserName string `db:"username"`
	Password string `db:"password"`
	Status   int64  `db:"status"`
}

type File struct {
	ID       int64  `json:"id" db:"id"`
	FileName string `json:"file_name" db:"file_name"`
	FileSize int64  `json:"file_size" db:"file_size"`
	Checksum string `json:"checksum" db:"file_checksum"`
	FileAddr string `json:"file_addr" db:"file_addr"`
	UserId   int64  `json:"user_id" db:"user_id"`
	Status   int64  `json:"status" db:"status"` //0表示已备份，1表示待备份，2表示待拉取
}
