package mongodb

import (
	"QQ-BOT/bot/db"
	"QQ-BOT/bot/model"
	"context"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const (
	MongoAdminAccountCollection      = "admin-account"
	MongoAssociationGroupCollection  = "association-group"
	MongoSchoolWallMessageCollection = "school-wall-message"
	MongoAuditRecordCollection       = "audit-record"
)

func init() {
	log.Println("mongodb init")
	db.RegisterDevice("mongodb", func() db.Database {
		data := viper.GetStringMap("database")
		mongo := data["mongodb"].(map[string]interface{})
		if mongo["enable"].(bool) {
			return &MongoDB{
				uri:   mongo["uri"].(string),
				db:    mongo["database"].(string),
				mongo: nil,
			}
		}
		return nil
	})
}

func (m *MongoDB) NextID(collectionName string) int64 {
	filter := bson.D{{"_id", collectionName}}
	update := bson.D{{"$inc", bson.D{{"value", 1}}}}
	var doc model.NextID
	err := m.mongo.Collection("counters").FindOneAndUpdate(context.TODO(), filter, update).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			m.mongo.Collection("counters").InsertOne(context.TODO(), bson.D{{"_id", collectionName}, {"value", 1}})
			return m.NextID(collectionName)
		}
		panic("nextID 错误")
	}
	return doc.Value
}

func (m *MongoDB) Find(collection string, filter interface{}, result interface{}) error {
	cur, err := m.mongo.Collection(collection).Find(context.TODO(), filter)
	if err != nil {
		return err
	}
	err = cur.All(context.TODO(), result)
	return err
}

func (m MongoDB) FindOne(collection string, filter interface{}, result interface{}) error {
	return m.mongo.Collection(collection).FindOne(context.TODO(), filter).Decode(result)
}

func (m *MongoDB) InsertOne(collection string, document interface{}) (err error) {
	_, err = m.mongo.Collection(collection).InsertOne(context.TODO(), document)
	return
}

func (m MongoDB) DeleteOne(collection string, filter interface{}) error {
	_, err := m.mongo.Collection(collection).DeleteOne(context.TODO(), filter)
	return err
}

type MongoDB struct {
	uri   string
	db    string
	mongo *mongo.Database
}

func (m *MongoDB) GetAuditRecord(mID int64) (record model.AuditRecord, err error) {
	err = m.FindOne(MongoAuditRecordCollection, bson.D{{"msg-id", mID}}, &record)
	return
}

func (m *MongoDB) AuditRecordIsExist(msgID int64) bool {
	var res []model.AuditRecord
	m.Find(MongoAuditRecordCollection, bson.D{{"msg-id", msgID}}, &res)
	if len(res) > 0 {
		return true
	}
	return false
}

func (m *MongoDB) AddAuditRecord(msgID int64, status string, operatorID int64, qqMID []int64) error {
	return m.InsertOne(MongoAuditRecordCollection, bson.D{{"msg-id", msgID}, {"status", status}, {"operatorID", operatorID}, {"qm-id", qqMID}})
}

func (m *MongoDB) GetSchoolWallMessage(msgID int64) (msg model.WallMessage, err error) {
	err = m.FindOne(MongoSchoolWallMessageCollection, bson.D{{"msg-id", msgID}}, &msg)
	return
}

func (m *MongoDB) AddSchoolWallMessage(msg string, school string) (msgID int64, err error) {
	msgID = m.NextID(MongoSchoolWallMessageCollection)
	err = m.InsertOne(MongoSchoolWallMessageCollection, bson.D{{"msg-id", msgID}, {"msg", msg}, {"school", school}})
	return
}

func (m *MongoDB) DeleteAdminAccount(account int64) error {
	return m.DeleteOne(MongoAdminAccountCollection, bson.D{{"account", account}})
}

func (m *MongoDB) AddAssociationGroup(groupID int64, operatorAccount int64, school string) error {
	err := m.InsertOne(MongoAssociationGroupCollection, bson.D{{"group-id", groupID}, {"operator-account", operatorAccount}, {"school", school}})
	return err
}

func (m *MongoDB) RemoveAssociationGroup(groupID int64) error {
	return m.DeleteOne(MongoAssociationGroupCollection, bson.D{{"group-id", groupID}})
}

func (m *MongoDB) GetAssociationGroupList() (data []model.AssociationGroup, err error) {
	err = m.Find(MongoAssociationGroupCollection, bson.D{}, &data)
	return
}

func (m *MongoDB) GetAdminAccountList() (data []model.AdminAccount, err error) {
	err = m.Find(MongoAdminAccountCollection, bson.D{}, &data)
	return
}

func (m *MongoDB) AddAdminAccount(account int64, parentAccount int64) error {
	_, err := m.mongo.Collection(MongoAdminAccountCollection).
		InsertOne(context.TODO(), bson.D{{"account", account}, {"parentAccount", parentAccount}})
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoDB) Open() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(m.uri))
	if err != nil {
		errors.Wrap(err, "open mongoDB connect error")
	}
	m.mongo = client.Database(m.db)
}
