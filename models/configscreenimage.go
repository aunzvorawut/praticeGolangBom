package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type ConfigScreenImage struct {
	Id            int64         `orm:"pk;auto"`
	Version       int64         `orm:"default(0)"`
	DateCreated   time.Time     `orm:"auto_now_add;type(datetime)"`
	LastUpdated   time.Time     `orm:"auto_now;type(datetime)"`
	OrderPosition int64         `orm:"default(0)"`
	ImagePathFat  string        `orm:"null;size(255)"`
	ImagePathThin string        `orm:"null;size(255)"`
	ConfigScreen  *ConfigScreen `orm:"null;rel(fk)"`
}

func init() {
	orm.RegisterModel(new(ConfigScreenImage))
}

// AddConfigScreenImage insert a new ConfigScreenImage into database and returns
// last inserted Id on success.
func AddConfigScreenImage(m *ConfigScreenImage) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetConfigScreenImageById retrieves ConfigScreenImage by Id. Returns error if
// Id doesn't exist
func GetConfigScreenImageById(id int64) (v *ConfigScreenImage, err error) {
	o := orm.NewOrm()
	v = &ConfigScreenImage{Id: id}
	if err = o.QueryTable(new(ConfigScreenImage)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllConfigScreenImage retrieves all ConfigScreenImage matches certain condition. Returns empty list if
// no records exist
func GetAllConfigScreenImage(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(ConfigScreenImage))
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

	var l []ConfigScreenImage
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

// UpdateConfigScreenImage updates ConfigScreenImage by Id and returns error if
// the record to be updated doesn't exist
func UpdateConfigScreenImageById(m *ConfigScreenImage) (err error) {
	o := orm.NewOrm()
	v := ConfigScreenImage{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteConfigScreenImage deletes ConfigScreenImage by Id and returns error if
// the record to be deleted doesn't exist
func DeleteConfigScreenImage(id int64) (err error) {
	o := orm.NewOrm()
	v := ConfigScreenImage{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&ConfigScreenImage{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetAllConfigScreenImageByconfigScreenObj(max, offset int64, configScreenObj *ConfigScreen) (ml []*ConfigScreenImage) {

	o := orm.NewOrm()
	qs := o.QueryTable(new(ConfigScreenImage))

	if _, err := qs.Limit(max, offset).Filter("ConfigScreen", configScreenObj).OrderBy("OrderPosition").RelatedSel().All(&ml); err == nil {
		return ml
	}

	return nil

}
