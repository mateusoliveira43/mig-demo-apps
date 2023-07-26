// MIT License

// Copyright (c) 2020 Mohamad Fadhil

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package database

import (
	"context"
	"errors"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sirupsen/logrus"
)

type ToDoTable struct {
	ToDo *mongo.Collection
}

type TodoItemModel struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	Description string
	Completed   bool
}

func (table ToDoTable) CreateItem(description string) string {
	logrus.WithFields(logrus.Fields{"description": description}).Info("Add new TodoItem. Saving to database.")
	todo := &TodoItemModel{Description: description, Completed: false}
	result, err := table.ToDo.InsertOne(context.TODO(), todo)
	if err != nil {
		// TODO gracefully fail?
		logrus.Fatal(err)
	}
	id := result.InsertedID.(primitive.ObjectID).Hex()
	logrus.Infof("inserted document with ID %v\n", id)
	return id
}

func (table ToDoTable) UpdateItem(id string, completed bool) (error, string) {
	if !table.GetItemByID(id) {
		return errors.New(""), "Record Not Found"
	}
	logrus.WithFields(logrus.Fields{"_id": id, "Completed": completed}).Info("Updating TodoItem")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		// TODO gracefully fail?
		panic(err)
	}
	filter := bson.M{"_id": objID}
	_, err = table.ToDo.UpdateOne(
		context.TODO(),
		filter,
		bson.D{
			{"$set", bson.D{{"completed", completed}}},
		},
	)
	if err != nil {
		return err, err.Error()
	}
	return nil, ""
}

func (table ToDoTable) DeleteItem(id string) (error, string) {
	if !table.GetItemByID(id) {
		return errors.New(""), "Record Not Found"
	}
	logrus.WithFields(logrus.Fields{"_id": id}).Info("Deleting TodoItem")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		// TODO gracefully fail?
		panic(err)
	}
	filter := bson.M{"_id": objID}
	opts := options.Delete().SetCollation(&options.Collation{
		Locale:    "en_US",
		Strength:  1,
		CaseLevel: false,
	})
	res, err := table.ToDo.DeleteOne(context.TODO(), filter, opts)
	if err != nil {
		// TODO gracefully fail?
		logrus.Fatal(err)
	}
	logrus.Infof("deleted %v documents\n", res.DeletedCount)
	return nil, ""
}

func (table ToDoTable) GetItemByID(Id string) bool {
	objID, err := primitive.ObjectIDFromHex(Id)
	if err != nil {
		// TODO gracefully fail?
		// TODO log.Fatal?
		panic(err)
	}
	filter := bson.M{"_id": objID}
	var result TodoItemModel
	err = table.ToDo.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		logrus.Error("ID NOT found")
		return false
	}
	logrus.Infof("%+v\n", result)
	return true
}

func (table ToDoTable) GetTodoItems(completed bool) []*TodoItemModel {
	findOptions := options.Find()
	findOptions.SetLimit(50)

	var results []*TodoItemModel
	filter := bson.M{"completed": completed}
	//TodoItems := db.Where("completed = ?", completed).Find(&todos).Value
	//return TodoItems
	cur, err := table.ToDo.Find(context.TODO(), filter, findOptions)
	if err != nil {
		logrus.Fatal(err)
	}

	// Iterate through the cursor
	for cur.Next(context.TODO()) {
		var elem TodoItemModel
		err := cur.Decode(&elem)
		if err != nil {
			logrus.Fatal(err)
		}

		results = append(results, &elem)
	}
	return results

}

func (table ToDoTable) PrePopulate() {
	// check to see if the db is prepopulated
	filter := bson.D{{"description", "time"}}
	var result TodoItemModel
	err := table.ToDo.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		logrus.Info("Prepopulate the db")
		prepop := TodoItemModel{Description: "prepopulate the db", Completed: true}
		donuts := TodoItemModel{Description: "time", Completed: false}
		both_prepop := []interface{}{prepop, donuts}

		insertManyResult, err := table.ToDo.InsertMany(context.TODO(), both_prepop)
		if err != nil {
			logrus.Fatal(err)
		}
		logrus.Info("Inserted multiple prepopulate documents: ", insertManyResult.InsertedIDs)
		// } else {
		// 	logrus.Infof("%+v\n", result)
	}
}

func GetToDoTable() *mongo.Collection {
	value, isSet := os.LookupEnv("ME_CONFIG_MONGODB_URL")
	var clientOptions *options.ClientOptions
	if isSet {
		clientOptions = options.Client().ApplyURI(value)
	}
	clientOptions = options.Client().ApplyURI("mongodb://changeme:changeme@localhost:27017")

	// Connect to MongoDB
	databaseConnection, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logrus.Fatal(err)
	}

	// Check the connection
	err = databaseConnection.Ping(context.TODO(), nil)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Connected to MongoDB!")
	// collection
	return databaseConnection.Database("todolist").Collection("TodoItemModel")
}
