package managers

import (
	"encoding/json"
	"net/http"
	d "restapi/dtos"
	m "restapi/models"
	u "restapi/utils"

	c "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// FindIntervalByID get an interval by ID
func FindIntervalByID(w http.ResponseWriter, r *http.Request) {
	if u.CheckMethodPath("GET", u.CheckRoutes["intervals"], w, r) {
		return
	}
	interval := []m.Interval{}
	vars := mux.Vars(r)
	if u.CheckErrorCode(u.GetConnection().Preload("Commentinterval").Preload("Commentinterval.User").Preload("Tags").Where("is_active = ?", true).Find(&interval, vars["id"]).Error, w) {
		return
	}
	u.Respond(w, interval)
}

// FindIntervalByAnnotationID get an interval by annotation ID
func FindIntervalByAnnotationID(w http.ResponseWriter, r *http.Request) {
	if u.CheckMethodPath("GET", u.CheckRoutes["intervalsannotations"], w, r) {
		return
	}
	interval := []m.Interval{}
	vars := mux.Vars(r)
	if u.CheckErrorCode(u.GetConnection().Preload("Commentinterval").Preload("Commentinterval.User").Preload("Tags").Where("is_active = ?", true).Where("annotation_id = ?", vars["id"]).Find(&interval).Error, w) {
		return
	}
	u.Respond(w, interval)
}

// CreateInterval create an interval
func CreateInterval(w http.ResponseWriter, r *http.Request) {
	if u.CheckMethodPath("POST", u.CheckRoutes["intervals"], w, r) {
		return
	}
	var i d.Interval
	db := u.GetConnection()
	annotation := m.Annotation{}
	contextUser := c.Get(r, "user").(*m.User)
	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil || i.TimeStart == nil || i.TimeEnd == nil || i.AnnotationID == nil && (*i.TimeStart > *i.TimeEnd) {
		http.Error(w, "Bad args", 400)
		return
	}

	if u.CheckErrorCode(db.Preload("Status").Preload("Status.EnumStatus").Preload("Status.User").Preload("Tags").First(&annotation, *i.AnnotationID).Error, w) {
		return
	}
	annotation.LastStatus, _ = annotation.GetLastAndFirstStatus()
	if annotation.LastStatus.EnumStatus.ID == 2 {
		if err := changeStatusEditDate(db, w, 3, &contextUser.ID, *i.AnnotationID); err != nil {
			if err == gorm.ErrRecordNotFound {
				http.Error(w, err.Error(), 404)
				return
			}
			if err == gorm.ErrInvalidSQL {
				http.Error(w, err.Error(), 400)
				return
			}
			http.Error(w, err.Error(), 500)
			return
		}
	}

	c := m.Interval{TimeStart: *i.TimeStart, TimeEnd: *i.TimeEnd, AnnotationID: *i.AnnotationID, IsActive: true}
	if u.CheckErrorCode(db.Create(&c).Error, w) {
		return
	}
	u.RespondCreate(w, c)
}

// RemoveIntervalByID Remove an interval by his ID
func RemoveIntervalByID(w http.ResponseWriter, r *http.Request) {
	if u.CheckMethodPath("DELETE", u.CheckRoutes["interval"], w, r) {
		return
	}
	interval := m.Interval{}
	v := mux.Vars(r)
	if len(v) != 1 || !u.IsStringInt(v["id"]) {
		http.Error(w, "Bad args", 204)
		return
	}
	db := u.GetConnection()
	if u.CheckErrorCode(db.First(&interval, v["id"]).Delete(&interval).Error, w) {
		return
	}
}

// AddTagsOnInterval create a tag on an interval
func AddTagsOnInterval(w http.ResponseWriter, r *http.Request) {
	if u.CheckMethodPath("POST", u.CheckRoutes["intervalstags"], w, r) {
		return
	}
	var i d.IntervalTagsPayload
	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil || i.Tags == nil || len(i.Tags) == 0 || i.IntervalID == nil {
		http.Error(w, "Bad request (client)", 204)
		return
	}
	tags := []m.Tag{}
	db := u.GetConnection()
	if u.CheckErrorCode(db.Find(&tags, i.Tags).Error, w) {
		return
	}
	if len(tags) != len(i.Tags) {
		http.Error(w, "Bad request (client)", 204)
		return
	}
	interval := m.Interval{}
	if u.CheckErrorCode(db.Find(&interval, *i.IntervalID).Error, w) {
		return
	}
	db.Model(&interval).Association("Tags").Replace(tags)
	u.RespondCreate(w, &tags)
}
