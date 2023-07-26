package main

import (
	"github.com/mig-demo-apps/apps/todolist-mongo-go/pkg/database"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	toDoTable  database.ToDoTable
	incomplete bool
)

var add = &cobra.Command{
	Use:    "add DESCRIPTION",
	Short:  "Add To Do item",
	Args:   cobra.ExactArgs(1),
	PreRun: connectToDatabase,
	Run: func(cmd *cobra.Command, args []string) {
		toDoTable.CreateItem(args[0])
	},
}

var update = &cobra.Command{
	Use:     "update ID",
	Short:   "Update To Do item to Complete or incomplete",
	Args:    cobra.ExactArgs(1),
	Example: "go run pkg/cli/cli.go update 64c11bd1da6b431c66c28a88",
	PreRun:  connectToDatabase,
	Run: func(cmd *cobra.Command, args []string) {
		toDoTable.UpdateItem(args[0], !incomplete)
	},
}

var delete = &cobra.Command{
	Use:     "delete ID",
	Short:   "Delete To Do item",
	Args:    cobra.ExactArgs(1),
	Example: "go run pkg/cli/cli.go delete 64c11bd1da6b431c66c28a88",
	PreRun:  connectToDatabase,
	Run: func(cmd *cobra.Command, args []string) {
		toDoTable.DeleteItem(args[0])
	},
}

var list = &cobra.Command{
	Use:     "list",
	Short:   "List To Do items",
	Args:    cobra.NoArgs,
	Example: "go run pkg/cli/cli.go list",
	PreRun:  connectToDatabase,
	Run: func(cmd *cobra.Command, args []string) {
		results := toDoTable.GetTodoItems(!incomplete)
		var done string
		if incomplete {
			done = "⬛ "
		} else {
			done = "✅ "
		}
		for _, result := range results {
			logrus.Info(done, result.Id.Hex(), " : ", result.Description)
		}
	},
}

var cli = &cobra.Command{
	Use:  "cli.go",
	Long: "Manage To Do items",
}

func connectToDatabase(cmd *cobra.Command, args []string) {
	toDoTable = database.ToDoTable{ToDo: database.GetToDoTable()}
	toDoTable.PrePopulate()
}

func init() {
	cli.AddCommand(add)
	cli.AddCommand(update)
	update.Flags().BoolVarP(&incomplete, "incomplete", "i", false, "Update to incomplete")
	cli.AddCommand(delete)
	cli.AddCommand(list)
	list.Flags().BoolVarP(&incomplete, "incomplete", "i", false, "List incomplete To Do items")

	cli.SetHelpCommand(&cobra.Command{Hidden: true})

	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetReportCaller(true)
}

func main() {
	err := cli.Execute()
	if err != nil {
		logrus.Fatal(err)
	}
}
