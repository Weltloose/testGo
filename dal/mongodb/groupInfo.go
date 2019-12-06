package mongodb

import (
	"fmt"
	"time"

	"github.com/Weltloose/testGo/dal/redis"
	"gopkg.in/mgo.v2/bson"
)

func CreateGroup(username string) int {
	id := redis.AddMaxGroupID()
	if id == -1 {
		return -1
	}
	group := GroupInfo{
		GroupID: id,
		Users: []string{
			username,
		},
	}
	err := groupInfoC.Insert(&group)
	if err != nil {
		fmt.Println("insert error ", err)
		return -1
	}
	if !userJoinGroup(username, id) {
		return -1
	}
	return id
}

func JoinGroup(username string, groupID int) bool {
	err := groupInfoC.Update(bson.M{"groupid": groupID}, bson.M{"$push": bson.M{"users": username}})
	if err != nil {
		fmt.Println("update error ", err)
		return false
	}
	if !userJoinGroup(username, groupID) {
		return false
	}
	return true
}

func LeaveGroup(username string, groupID int) bool {
	err := groupInfoC.Update(bson.M{"groupid": groupID}, bson.M{"$pull": bson.M{"users": username}})
	if err != nil {
		fmt.Println("update error ", err)
		return false
	}
	if !userLeaveGroup(username, groupID) {
		return false
	}
	return true
}

func GetEvents(groupID int) []Event {
	var res GroupInfo
	err := groupInfoC.Find(bson.M{"groupid": groupID}).One(&res)
	if err != nil {
		fmt.Println("get groupInfo error ", err)
		return []Event{}
	}
	return res.Events
}

func AddEvent(groupID int, fromTime, toTime, Content string) {
	var res GroupInfo
	err := groupInfoC.Find(bson.M{"groupid": groupID}).One(&res)
	if err != nil {
		fmt.Println("get groupInfo error ", err)
		return
	}
	evList := res.Events
	eid := redis.GetNewEventID()
	newEvent := Event{
		EventID:  eid,
		TimeFrom: fromTime,
		TimeTo:   toTime,
		Content:  Content,
	}
	var opt []Event
	ok := false
	newEventTimeFrom, _ := time.ParseInLocation("2006-01-02 15:04:05", fromTime, time.Local)
	for i, val := range evList {
		valTimefrom, _ := time.ParseInLocation("2006-01-02 15:04:05", val.TimeFrom, time.Local)
		if valTimefrom.After(newEventTimeFrom) {
			opt = append(opt, evList[:i]...)
			opt = append(opt, newEvent)
			opt = append(opt, evList[i:]...)
			ok = true
			break
		}
	}
	if !ok {
		opt = append(opt, evList...)
		opt = append(opt, newEvent)
	}
	groupInfoC.Update(bson.M{"groupid": groupID}, bson.M{"$set": bson.M{"events": opt}})
	wsret := WSopt{
		GroupID: groupID,
		Events:  opt,
	}
	MainHub.Broadcast <- wsret
}

func EditEvent(groupID, eventID int, fromTime, toTime, Content string) {
	var res GroupInfo
	err := groupInfoC.Find(bson.M{"groupid": groupID}).One(&res)
	if err != nil {
		fmt.Println("get groupInfo error ", err)
		return
	}
	evList := res.Events
	for i, val := range evList {
		if val.EventID == eventID {
			evList[i].TimeFrom = fromTime
			evList[i].TimeTo = toTime
			evList[i].Content = Content
		}
	}
	// 修改后按TimeFrom重新排序好
	for i, _ := range evList {
		evListITimeFrom, _ := time.ParseInLocation("2006-01-02 15:04:05", evList[i].TimeFrom, time.Local)
		for j := i + 1; j < len(evList); j++ {
			evListJTimeFrom, _ := time.ParseInLocation("2006-01-02 15:04:05", evList[j].TimeFrom, time.Local)
			if evListITimeFrom.After(evListJTimeFrom) {
				tmp := evList[i]
				evList[i] = evList[j]
				evList[j] = tmp
			}
		}
	}
	groupInfoC.Update(bson.M{"groupid": groupID}, bson.M{"$set": bson.M{"events": evList}})
	wsret := WSopt{
		GroupID: groupID,
		Events:  evList,
	}
	MainHub.Broadcast <- wsret
}

func DeleteEvent(groupID, eventID int) {
	var res GroupInfo
	err := groupInfoC.Find(bson.M{"groupid": groupID}).One(&res)
	if err != nil {
		fmt.Println("get groupInfo error ", err)
		return
	}
	evList := res.Events
	for i, val := range evList {
		if val.EventID == eventID {
			optList := append(evList[:i], evList[i+1:]...)
			groupInfoC.Update(bson.M{"groupid": groupID}, bson.M{"$set": bson.M{"events": optList}})
			wsret := WSopt{
				GroupID: groupID,
				Events:  optList,
			}
			MainHub.Broadcast <- wsret
			return
		}
	}
}
