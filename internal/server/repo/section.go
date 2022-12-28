package serverRepo

import (
	"strconv"
	"strings"

	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/gen"
	"github.com/easysoft/zendata/internal/pkg/model"
	"gorm.io/gorm"
)

type SectionRepo struct {
	DB *gorm.DB `inject:""`
}

func (r *SectionRepo) List(ownerId uint, ownerType string) (sections []*model.ZdSection, err error) {
	err = r.DB.Where("ownerID=? AND ownerType=?", ownerId, ownerType).Find(&sections).Error
	return
}

func (r *SectionRepo) Get(id uint, ownerType string) (section model.ZdSection, err error) {
	err = r.DB.Where("id=? AND ownerType=?", id, ownerType).First(&section).Error
	return
}

func (r *SectionRepo) Create(section *model.ZdSection) (err error) {
	err = r.DB.Create(&section).Error
	return
}

func (r *SectionRepo) Update(section *model.ZdSection) (err error) {
	err = r.DB.Save(&section).Error
	return
}

func (r *SectionRepo) Remove(id uint, ownerType string) (err error) {
	err = r.DB.Where("id=? AND ownerType=?", id, ownerType).Delete(&model.ZdSection{}).Error
	return
}

func (r *SectionRepo) SaveFieldSectionToDB(rangeSection string, ord int, fieldID uint, ownerType string) {
	descStr, stepStr, count, countTag := gen.ParseRangeSection(rangeSection)
	typ, desc := gen.ParseRangeSectionDesc(descStr)

	if typ == "literal" && desc[:1] == string(consts.LeftBrackets) &&
		desc[len(desc)-1:] == string(consts.RightBrackets) {

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

	start := ""
	end := ""
	if typ == "interval" {
		arr := strings.Split(desc, "-")
		start = arr[0]
		if len(arr) > 1 {
			end = arr[1]
		}
	}

	section := model.ZdSection{OwnerType: ownerType, OwnerID: fieldID, Type: typ,
		Value: desc, Start: start, End: end, Ord: ord,
		Step: step, Repeat: countStr, RepeatTag: countTag, Rand: rand}

	r.Create(&section)
}

func NewSectionRepo(db *gorm.DB) *SectionRepo {
	return &SectionRepo{DB: db}
}
