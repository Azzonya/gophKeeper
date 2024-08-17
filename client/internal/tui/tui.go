// Package tui provides a text-based user interface (TUI) for the GophKeeper client,
// allowing users to register, log in, and manage their data items through a simple
// interactive interface using tview.
package tui

import (
	"context"
	"google.golang.org/grpc/metadata"
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"gophKeeper/client/internal/client"
	proto "gophKeeper/pkg/proto/gophkeeper"
)

// TUI represents the text-based user interface for the GophKeeper client,
// handling the display and interaction logic for user registration, login,
// and data item management.
type TUI struct {
	client *client.GophKeeperClient
	app    *tview.Application
}

// NewTUI creates a new TUI instance with the given gRPC client, initializing
// the application and setting up the user interface.
func NewTUI(client *client.GophKeeperClient) *TUI {
	return &TUI{
		client: client,
		app:    tview.NewApplication(),
	}
}

// Run starts the TUI application, displaying the main form with options for
// user registration, login, and quitting the application.
func (t *TUI) Run() error {
	form := tview.NewForm()

	form.AddButton("Register", t.register).
		AddButton("Login", t.login).
		AddButton("Quit", func() {
			t.app.Stop()
		})

	t.app.SetRoot(form, true)
	return t.app.Run()
}

// register handles the user registration process, displaying a form to input
// a username and password and sending the registration request to the server.
func (t *TUI) register() {
	form := tview.NewForm()

	form.
		AddInputField("Username", "", 20, nil, nil).
		AddInputField("Password", "", 20, nil, nil).
		AddButton("Register", func() {
			username := form.GetFormItemByLabel("Username").(*tview.InputField).GetText()
			password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			resp, err := t.client.Register(ctx, &proto.RegisterRequest{
				Username: username,
				Password: password,
			})
			if err != nil {
				log.Printf("Register failed: %v", err)
				t.showMessage("Register failed. Press Enter to go back.", t.restart)
				return
			}

			log.Printf("Register successful, message: %s", resp.Message)
			t.showMessage("Register successful. Press Enter to go back.", t.restart)
		}).
		AddButton("Cancel", func() {
			t.restart()
		})

	t.app.SetRoot(form, true).SetFocus(form)
}

// login handles the user login process, displaying a form to input a username
// and password and sending the login request to the server.
func (t *TUI) login() {
	form := tview.NewForm()
	form.
		AddInputField("Username", "", 20, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil).
		AddButton("Login", func() {
			username := form.GetFormItemByLabel("Username").(*tview.InputField).GetText()
			password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs())

			resp, err := t.client.Login(ctx, &proto.LoginRequest{
				Username: username,
				Password: password,
			})
			if err != nil {
				log.Printf("Login failed: %v", err)
				t.showMessage("Login failed. Press Enter to go back.", t.restart)
				return
			}

			t.client.BearerToken = resp.Token

			log.Printf("Login successful")
			t.showMainMenu()
		}).
		AddButton("Cancel", func() {
			t.restart()
		})

	t.app.SetRoot(form, true).SetFocus(form)
}

// showMainMenu displays the main menu with options for creating, getting,
// updating, and deleting data items, as well as quitting the application.
func (t *TUI) showMainMenu() {
	menu := tview.NewList().
		AddItem("Create Data", "Create new data", 'c', t.createData).
		AddItem("Get Data", "Get existing data", 'g', t.getData).
		AddItem("Update Data", "Update existing data", 'u', t.updateData).
		AddItem("Delete Data", "Delete existing data", 'd', t.deleteData).
		AddItem("Quit", "Press to exit", 'q', func() {
			t.app.Stop()
		})

	t.app.SetRoot(menu, true).SetFocus(menu)
}

// createData displays a form for creating a new data item, allowing the user
// to input the type, data, and metadata, and sending the create request to the server.
func (t *TUI) createData() {
	form := tview.NewForm()

	form.
		AddDropDown("Type", []string{"binary", "text", "credentials", "bank card"}, 0, nil).
		AddInputField("Data", "", 20, nil, nil).
		AddInputField("Meta", "", 20, nil, nil).
		AddButton("Submit", func() {
			_, typeField := form.GetFormItemByLabel("Type").(*tview.DropDown).GetCurrentOption()
			dataField := form.GetFormItemByLabel("Data").(*tview.InputField).GetText()
			metaField := form.GetFormItemByLabel("Meta").(*tview.InputField).GetText()

			req := &proto.CreateDataRequest{
				Data: &proto.DataItem{
					Type: typeField,
					Data: []byte(dataField),
					Meta: metaField,
				},
			}

			md := metadata.Pairs(
				"token", "Bearer "+t.client.BearerToken,
			)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
			defer cancel()

			ctx = metadata.NewOutgoingContext(ctx, md)

			resp, err := t.client.CreateData(ctx, req)
			if err != nil {
				log.Printf("failed to create data: %v", err)
				t.showMessage("Failed to create data. Press Enter to go back.", t.showMainMenu)
				return
			}
			log.Printf("CreateData response: %s", resp.Message)
			t.showMessage("Data created successfully. Press Enter to go back.", t.showMainMenu)
		}).
		AddButton("Cancel", func() {
			t.showMainMenu()
		})

	t.app.SetRoot(form, true).SetFocus(form)
}

// getData displays a form for retrieving an existing data item, allowing the user
// to input the ID and type of the data item and sending the get request to the server.
func (t *TUI) getData() {
	form := tview.NewForm()
	form.
		AddInputField("ID", "", 20, nil, nil).
		AddInputField("Type", "", 20, nil, nil).
		AddButton("Submit", func() {
			idField := form.GetFormItemByLabel("ID").(*tview.InputField).GetText()
			typeField := form.GetFormItemByLabel("Type").(*tview.InputField).GetText()
			URLField := form.GetFormItemByLabel("Type").(*tview.InputField).GetText()

			req := &proto.GetDataRequest{
				Id:   idField,
				Type: typeField,
				URL:  URLField,
			}

			md := metadata.Pairs(
				"token", "Bearer "+t.client.BearerToken,
			)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			ctx = metadata.NewOutgoingContext(ctx, md)

			resp, err := t.client.GetData(ctx, req)
			if err != nil {
				log.Printf("failed to get data: %v", err)
				t.showMessage("Failed to get data. Press Enter to go back.", t.showMainMenu)
				return
			}
			if len(resp.Data) > 0 {
				log.Printf("GetData response: %s", string(resp.Data[0].Data))
				t.showMessage("Data retrieved successfully. Press Enter to go back.", t.showMainMenu)
			} else {
				log.Printf("GetData response: no data found")
				t.showMessage("No data found. Press Enter to go back.", t.showMainMenu)
			}
		}).
		AddButton("Cancel", func() {
			t.showMainMenu()
		})

	t.app.SetRoot(form, true).SetFocus(form)
}

// updateData displays a form for updating an existing data item, allowing the user
// to input the ID, type, data, and metadata, and sending the update request to the server.
func (t *TUI) updateData() {
	form := tview.NewForm()
	form.
		AddInputField("ID", "", 20, nil, nil).
		AddInputField("Type", "", 20, nil, nil).
		AddInputField("Data", "", 20, nil, nil).
		AddInputField("Meta", "", 20, nil, nil).
		AddButton("Submit", func() {
			idField := form.GetFormItemByLabel("ID").(*tview.InputField).GetText()
			typeField := form.GetFormItemByLabel("ID").(*tview.InputField).GetText()
			dataField := form.GetFormItemByLabel("ID").(*tview.InputField).GetText()
			metaField := form.GetFormItemByLabel("ID").(*tview.InputField).GetText()

			req := &proto.UpdateDataRequest{
				Data: &proto.DataItem{
					Id:   idField,
					Type: typeField,
					Data: []byte(dataField),
					Meta: metaField,
				},
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			resp, err := t.client.UpdateData(ctx, req)
			if err != nil {
				log.Printf("failed to get data: %v", err)
				t.showMessage("Failed to get data. Press Enter to go back.", t.showMainMenu)
				return
			}
			if len(resp.Message) > 0 {
				log.Printf("UpdateData response: %s", resp.Message)
				t.showMessage("Data updated successfully. Press Enter to go back.", t.showMainMenu)
			} else {
				log.Printf("UpdateData response: no data found")
				t.showMessage("No data found. Press Enter to go back.", t.showMainMenu)
			}
		}).
		AddButton("Cancel", func() {
			t.showMainMenu()
		})

	t.app.SetRoot(form, true).SetFocus(form)
}

// deleteData displays a form for deleting a data item, allowing the user
// to input the ID and sending the delete request to the server.
func (t *TUI) deleteData() {
	form := tview.NewForm()
	form.
		AddInputField("ID", "", 20, nil, nil).
		AddButton("Submit", func() {
			idField := form.GetFormItemByLabel("ID").(*tview.InputField).GetText()

			req := &proto.DeleteDataRequest{
				Id: idField,
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			resp, err := t.client.DeleteData(ctx, req)
			if err != nil {
				log.Printf("failed to get data: %v", err)
				t.showMessage("Failed to get data. Press Enter to go back.", t.showMainMenu)
				return
			}
			if len(resp.Message) > 0 {
				log.Printf("DeleteData response: %s", resp.Message)
				t.showMessage("Data deleted successfully. Press Enter to go back.", t.showMainMenu)
			} else {
				log.Printf("DeleteData response: no data found")
				t.showMessage("No data found. Press Enter to go back.", t.showMainMenu)
			}
		}).
		AddButton("Cancel", func() {
			t.showMainMenu()
		})

	t.app.SetRoot(form, true).SetFocus(form)
}

// showMessage displays a message to the user with a prompt to press Enter to continue,
// returning to a specified function after the message is acknowledged.
func (t *TUI) showMessage(message string, doneFunc func()) {
	textView := tview.NewTextView().
		SetText(message).
		SetDoneFunc(func(key tcell.Key) {
			doneFunc()
		})
	t.app.SetRoot(textView, true).SetFocus(textView)
}

// restart restarts the TUI application, resetting the interface and returning
// to the initial state.
func (t *TUI) restart() {
	// Перезапуск приложения
	go func() {
		if err := t.Run(); err != nil {
			log.Fatalf("failed to run TUI: %v", err)
		}
	}()
}
