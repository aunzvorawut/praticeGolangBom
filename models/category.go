package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Category struct {
	Id              int64     `orm:"pk;auto"`
	Version         int64     `orm:"default(0)"`
	DateCreated     time.Time `orm:"auto_now_add;type(datetime)"`
	LastUpdated     time.Time `orm:"auto_now;type(datetime)"`
	CategoryName    string
	CategoryNameEng string
	Cover           string
	CoverImage      string
	OrderPosition   string
	IsEnabled       bool `orm:"null;default(true)"`
	IsFree          bool `orm:"null;default(false)"`
	IsRecordable    bool `orm:"null;default(false)"`
	CatCodeJob      string
	IsBox           bool `orm:"null;default(true)"`
	IsMobile        bool `orm:"null;default(true)"`
	CategoryType    string
}

func init() {
	orm.RegisterModel(new(Category))
}

// AddCategory insert a new Category into database and returns
// last inserted Id on success.
func AddCategory(m *Category) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetCategoryById retrieves Category by Id. Returns error if
// Id doesn't exist
func GetCategoryById(id int64) (v *Category, err error) {
	o := orm.NewOrm()
	v = &Category{Id: id}
	if err = o.QueryTable(new(Category)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllCategory retrieves all Category matches certain condition. Returns empty list if
// no records exist
func GetAllCategory(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Category))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Category
	qs = qs.OrderBy(sortFields...).RelatedSel()
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateCategory updates Category by Id and returns error if
// the record to be updated doesn't exist
func UpdateCategoryById(m *Category) (err error) {
	o := orm.NewOrm()
	v := Category{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteCategory deletes Category by Id and returns error if
// the record to be deleted doesn't exist
func DeleteCategory(id int64) (err error) {
	o := orm.NewOrm()
	v := Category{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Category{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetCategoryFreeRandom() *Category {
	o := orm.NewOrm()
	categoryObj := new(Category)
	err := o.QueryTable(new(Category)).Filter("IsFree", true).RelatedSel().Limit(1).One(categoryObj)
	if err != nil {
		return nil
	}
	return categoryObj
}

func GetCategoryTrendingRandom() *Category {
	o := orm.NewOrm()
	categoryObj := new(Category)
	err := o.QueryTable(new(Category)).Filter("CategoryType", "Trending").RelatedSel().Limit(1).One(categoryObj)
	if err != nil {
		return nil
	}
	return categoryObj
}

func GetAllCategoryByClientType(max, offset int, clientTypeFunc string) (ml []*Category) {

	o := orm.NewOrm()
	qs := o.QueryTable(new(Category))
	qs = qs.Limit(max, offset).Filter("IsEnabled", true)

	if clientTypeFunc == "box" {
		qs = qs.Filter("IsBox", true)
	}

	if clientTypeFunc == "mobile" {
		qs = qs.Filter("IsMobile", true)
	}

	_, err := qs.OrderBy("OrderPosition").RelatedSel().All(&ml)

	if err == nil {
		return ml
	}

	return nil

}

func GetAllCategoryByisFree(max, offset int, isFree bool) (ml []*Category) {

	o := orm.NewOrm()
	qs := o.QueryTable(new(Category))

	_, err := qs.Limit(max, offset).Filter("IsFree", isFree).RelatedSel().All(&ml)

	if err == nil {
		return ml
	}

	return nil

}
