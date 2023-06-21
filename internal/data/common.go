package data

import (
	"errors"
	"fmt"
	"time"

	"alsritter.icu/rabbit-template/api/common"
	"alsritter.icu/rabbit-template/internal/pkg/timetools"

	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrModelNotBeNil         = errors.New("model not be nil")
	ErrIdLessThanZero        = errors.New("id not less than zero")
	ErrModelsNotBeEmptySlice = errors.New("models not be empty slice")
	ErrConditionNotBeNil     = errors.New("condition not be nil")
	ErrIdsNotBeEmpty         = errors.New("ids not be empty array")
	ErrOmnipotent            = errors.New("系统繁忙，请稍后再试")
	ErrNetwork               = errors.New("网络异常,请稍后再试")
	ErrDBFail                = errors.New("数据库异常")
	ErrExpire                = errors.New("账号已经过期请重新登录")
	ErrCatnotPay             = errors.New("用户未支付")
	ErrParam                 = errors.New("参数错误")
	ErrAMapDistrictNotFound  = errors.New("地图未找到对应的行政区域")
	ErrWxChatSDKInit         = errors.New("微信小程序SDK初始化失败")
)

type Model struct {
	ID        int32      `gorm:"primary_key;AUTO_INCREMENT"`
	CreatedAt *time.Time `json:"create_time" db:"create_time" gorm:"autoCreateTime"`
	CreatedBy string
	UpdatedAt *time.Time `json:"modify_time" db:"modify_time" gorm:"autoUpdateTime"`
	UpdatedBy string
	DeletedAt gorm.DeletedAt
	DeletedBy string
}

type ModelT interface {
	TableName() string
}

type CommonRepo[M ModelT] struct {
	db *gorm.DB
}

// DB 回调函数的选项
type DBOption[M ModelT] func(r CommonRepo[M])

// newCommonRepo 创建一个通用的 repo
func newCommonRepo(model ModelT, db *gorm.DB, options ...DBOption[ModelT]) CommonRepo[ModelT] {
	repo := CommonRepo[ModelT]{
		db: db.Table(model.TableName()),
	}
	for _, opt := range options {
		opt(repo)
	}
	return repo
}

func (r *CommonRepo[M]) initDB() *gorm.DB {
	db := r.db.Session(&gorm.Session{})
	return db
}

func (r *CommonRepo[M]) Create(model *M) error {
	db := r.initDB().Create(model)
	if err := db.Error; err != nil {
		return err
	}
	return nil
}

func (r *CommonRepo[M]) CreateBatch(models *[]M) error {
	db := r.initDB().Create(models)
	if err := db.Error; err != nil {
		return err
	}
	return nil
}

func (r *CommonRepo[M]) FirstOrCreate(condition *M) (*M, error) {
	var newModel M
	db := r.initDB().FirstOrCreate(&newModel, condition)
	if err := db.Error; err != nil {
		return nil, err
	}
	return &newModel, nil
}

func (r *CommonRepo[M]) Update(model *M, condition *M, fields ...string) error {
	db := r.initDB()
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	db = db.Where(condition).Updates(model)
	if err := db.Error; err != nil {
		return err
	}

	return nil
}

func (r *CommonRepo[M]) UpdateFiledById(model *M, id int32, fields ...string) error {
	db := r.initDB()
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	db = db.Where("id = ?", id).Updates(model)
	if err := db.Error; err != nil {
		return err
	}
	return nil
}

func (r *CommonRepo[M]) UpdateById(model *M, id int32) error {
	db := r.initDB().Where("id = ?", id).Updates(model)
	if err := db.Error; err != nil {
		return err
	}
	return nil
}

func (r *CommonRepo[M]) UpdateAllById(model *M, id int32) error {
	db := r.initDB().Select("*").Where("id = ?", id).Updates(model)
	if err := db.Error; err != nil {
		return err
	}
	return nil
}

// more see: https://gorm.io/zh_CN/docs/create.html#Upsert-%E5%8F%8A%E5%86%B2%E7%AA%81
func (r *CommonRepo[M]) Upsert(models *[]M) error {
	db := r.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(models)

	if err := db.Error; err != nil {
		return err
	}
	return nil
}

// 通过条件删除记录
func (r *CommonRepo[M]) Delete(condition *M, deleteUser ...string) error {
	if condition == nil {
		return ErrConditionNotBeNil
	}
	var deleteUser0 string
	if len(deleteUser) > 0 {
		deleteUser0 = deleteUser[0]
	}
	db := r.db.Where("deleted_at = ?", 0).Where(condition).Updates(map[string]interface{}{
		"deleted_at": time.Now(),
		"deleted_by": deleteUser0,
	})
	if err := db.Error; err != nil {
		return err
	}
	return nil
}

// 通过 IDs 删除记录
func (r *CommonRepo[M]) DeleteByIds(ids []int32, deleteUser ...string) error {
	if len(ids) <= 0 {
		return ErrIdsNotBeEmpty
	}
	var deleteUser0 string
	if len(deleteUser) > 0 {
		deleteUser0 = deleteUser[0]
	}

	db := r.db.Where(`deleted_at IS NULL AND id IN ?`, ids).Updates(map[string]interface{}{
		"deleted_at": time.Now(),
		"deleted_by": deleteUser0,
	})
	if err := db.Error; err != nil {
		return err
	}
	return nil
}

// 通过 ID 删除记录
func (r *CommonRepo[M]) DeleteById(id int32, deleteUser ...string) error {
	if id <= 0 {
		return ErrIdLessThanZero
	}
	err := r.DeleteByIds([]int32{id}, deleteUser...)
	if err != nil {
		return err
	}
	return nil
}

// 通过条件查询记录
func (r *CommonRepo[M]) Get(condition *M, fields ...string) (*M, error) {
	result, err := r.GetExist(condition, fields...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

// 通过条件查询数量
func (r *CommonRepo[M]) GetCount(condition *M) (int32, error) {
	var count int64
	if err := r.initDB().Model(condition).
		Where(condition).Where("deleted_at is null").
		Count(&count).Error; err != nil {
		return 0, err
	}
	return int32(count), nil
}

// 通过 ID 查询记录，如果不存在则会报错
func (r *CommonRepo[M]) GetExist(condition *M, fields ...string) (*M, error) {
	r.db = r.initDB()
	var record M
	db := r.db.Where(condition)
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	db = db.First(&record)
	if err := db.Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// 通过 ID 查询记录
func (r *CommonRepo[M]) GetById(id int32, fields ...string) (*M, error) {
	result, err := r.GetExistById(id, fields...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

// 通过 ID 查询记录(包括被软删除的数据)
func (r *CommonRepo[M]) GetExistUnscopedById(id int32, fields ...string) (*M, error) {
	if id <= 0 {
		return nil, ErrIdLessThanZero
	}
	var record M
	db := r.initDB()
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	db = db.Unscoped().First(&record, id)
	if err := db.Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// 通过 ID 查询记录，如果不存在则会报错
func (r *CommonRepo[M]) GetExistById(id int32, fields ...string) (*M, error) {
	if id <= 0 {
		return nil, ErrIdLessThanZero
	}
	var record M
	db := r.initDB()
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	db = db.First(&record, id)
	if err := db.Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// 通过条件查询是否存在记录
func (r *CommonRepo[M]) Exist(condition *M) (bool, error) {
	_, err := r.GetExist(condition)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

// 通过条件查询是否存在记录
func (r *CommonRepo[M]) ExistById(id int32) (bool, error) {
	_, err := r.GetExistById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

// 通过条件查询是否存在记录，排除指定 IDs
func (r *CommonRepo[M]) ExistExcludeIds(condition *M, ids []int32) (bool, error) {
	if len(ids) <= 0 {
		return false, ErrIdsNotBeEmpty
	}

	var record M
	db := r.initDB().Where("id NOT IN ?", ids).Where(condition)
	db = db.First(&record)
	if err := db.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

// 通过条件查询是否存在记录，排除指定 ID
func (r *CommonRepo[M]) ExistExcludeId(condition *M, id int32) (bool, error) {
	if id <= 0 {
		return false, ErrIdLessThanZero
	}
	return r.ExistExcludeIds(condition, []int32{id})
}

// 通过条件查询记录列表
func (r *CommonRepo[M]) List(condition *M, fields ...string) ([]*M, error) {
	var records []*M
	db := r.initDB().Where(condition)
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	db = db.Find(&records)
	if err := db.Error; err != nil {
		return records, err
	}
	return records, nil
}

// 通过条件查询记录列表（包含删除的记录）
func (r *CommonRepo[M]) ListAll(condition *M, fields ...string) ([]*M, error) {
	var records []*M
	db := r.initDB().Unscoped().Where(condition)
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	db = db.Find(&records)
	if err := db.Error; err != nil {
		return records, err
	}
	return records, nil
}

// 通过条件分页查询记录列表
func (r *CommonRepo[M]) Page(condition *M, pager *common.Pager, fields ...string) ([]*M, error) {
	records, err := r.PageOrder(condition, pager, "updated_at desc", fields...)
	if err != nil {
		return nil, err
	}
	return records, nil
}

// 通过条件分页查询记录列表（指定排序）
func (r *CommonRepo[M]) PageOrder(condition *M, pager *common.Pager, order string, fields ...string) ([]*M, error) {
	var records []*M
	db := r.initDB().Model(condition).Where(condition)
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	if len(order) > 0 {
		db = db.Order(order)
	}
	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return nil, err
	}

	pager.Total = int32(count)
	db = db.Offset(getOffset(pager)).Limit(int(pager.PageSize)).Find(&records)
	if err := db.Error; err != nil {
		return records, err
	}
	return records, nil
}

func getOffset(pager *common.Pager) int {
	if pager.Page <= 0 {
		return 0
	}
	return int((pager.Page - 1) * pager.PageSize)
}

func FillModel(model Model) Model {
	if model.CreatedAt == nil {
		model.CreatedAt = timetools.Now()
	}
	if model.UpdatedAt == nil {
		model.UpdatedAt = timetools.Now()
	}
	return model
}

func GenerateNumber() int32 {
	node, _ := snowflake.NewNode(16)
	return int32(node.Generate())
}

func GenerateTaderNo() string {
	node, _ := snowflake.NewNode(16)
	return fmt.Sprintf("%d", node.Generate())
}
