package mongodb

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

func UserNameExist(username string) bool {

	n, err := userInfoC.Find(bson.M{"username": username}).Count()
	if err != nil || n == 1 {
		return true
	}
	return false
}

func AddUser(username, password string) bool {
	user := UserInfo{
		Username: username,
		Password: password,
	}
	err := userInfoC.Insert(&user)
	if err != nil {
		fmt.Println("insert error ", err)
		return false
	}
	return true
}

func ValidateUser(username, password string) bool {
	n, err := userInfoC.Find(bson.M{"username": username, "password": password}).Count()
	if err != nil || n == 1 {
		return true
	}
	return false
}

func userJoinGroup(username string, groupID int) bool {
	err := userInfoC.Update(bson.M{"username": username}, bson.M{"$push": bson.M{"groups": groupID}})
	if err != nil {
		fmt.Println("update error ", err)
		return false
	}
	return true
}

func userLeaveGroup(username string, groupID int) bool {
	err := userInfoC.Update(bson.M{"username": username}, bson.M{"$pull": bson.M{"groups": groupID}})
	if err != nil {
		fmt.Println("update error ", err)
		return false
	}
	return true
}

func GetGroupList(username string) []int {
	var user UserInfo
	err := userInfoC.Find(bson.M{"username": username}).One(&user)
	if err != nil {
		return []int{}
	}
	return user.Groups
}

func GetEventList(username string) EventListOpt {
	var opt EventListOpt
	opt.Data = make(map[int][]Event)
	groups := GetGroupList(username)
	for _, val := range groups {
		opt.Data[val] = GetEvents(val)
	}
	return opt
}
