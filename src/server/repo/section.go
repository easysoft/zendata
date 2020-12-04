package serverRepo

import (
	"github.com/easysoft/zendata/src/gen"
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	"github.com/jinzhu/gorm"
	"strconv"
)

type SectionRepo struct {
	db *gorm.DB
}

func (r *SectionRepo) List(ownerId uint, ownerType string) (sections []*model.ZdSection, err error) {
	err = r.db.Where("ownerID=? AND ownerType=?", ownerId, ownerType).Find(&sections).Error
	return
}

func (r *SectionRepo) Get(id uint, ownerType string) (section model.ZdSection, err error) {
	err = r.db.Where("id=? AND ownerType=?", id, ownerType).First(&section).Error
	return
}

func (r *SectionRepo) Create(section *model.ZdSection) (err error) {
	err = r.db.Create(&section).Error
	return
}

func (r *SectionRepo) Update(section *model.ZdSection) (err error) {
	err = r.db.Save(&section).Error
	return
}

func (r *SectionRepo) Remove(id uint, ownerType string) (err error) {
	err = r.db.Where("id=? AND ownerType=?", id, ownerType).Delete(&model.ZdSection{}).Error
	return
}

func (r *SectionRepo) SaveFieldSectionToDB(rangeSection string, ord int, fieldID uint, ownerType string) {
	descStr, stepStr, count := gen.ParseRangeSection(rangeSection)
	typ, desc := gen.ParseRangeSectionDesc(descStr)

	if typ == "literal" && desc[:1] == string(constant.LeftBrackets) &&
		desc[len(desc)-1:] == string(constant.RightBrackets) {

		desc = "[" + desc[1:len(desc)-1] + "]"
		typ = "list"
	}

	countStr := strconv.Itoa(count)
	rand := false
	step := 1
	if stepStr == "r" {
		rand = true
	} else {
		step, _ = strconv.Atoi(stepStr)
	}
	section := model.ZdSection{OwnerType: ownerType, OwnerID: fieldID, Type: typ, Value: desc, Ord: ord,
		Step: step, Repeat: countStr, Rand: rand}

	r.Create(&section)
}

func NewSectionRepo(db *gorm.DB) *SectionRepo {
	return &SectionRepo{db: db}
}
