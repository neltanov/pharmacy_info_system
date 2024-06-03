package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type QueryRequest struct {
	Query  string                 `json:"query"`
	Params map[string]interface{} `json:"params"`
}

type QueryResult struct {
	Columns []string        `json:"columns"`
	Rows    [][]interface{} `json:"rows"`
}

type Order struct {
	ID             int    `json:"id"`
	CustomerID     int    `json:"customer_id"`
	ReceiptID      int    `json:"receipt_id"`
	OrderDate      string `json:"order_date"`
	ProductionDate string `json:"production_date"`
	Status         string `json:"status"`
}

type Customer struct {
	ID          int    `json:"id"`
	Surname     string `json:"surname"`
	Name        string `json:"name"`
	MiddleName  string `json:"middle_name"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
}

type Patient struct {
	ID         int    `json:"id"`
	Surname    string `json:"surname"`
	Name       string `json:"name"`
	MiddleName string `json:"middle_name"`
	Birthday   string `json:"birthday"`
}

type Doctor struct {
	ID         int    `json:"id"`
	Surname    string `json:"surname"`
	Name       string `json:"name"`
	MiddleName string `json:"middle_name"`
	Specialty  string `json:"specialty"`
}

type Medicine struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Receipt struct {
	ID        int `json:"id"`
	PatientID int `json:"patient_id"`
	DoctorID  int `json:"doctor_id"`
}

func sendQuery(query string, params map[string]interface{}) (QueryResult, error) {
	req := QueryRequest{
		Query:  query,
		Params: params,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return QueryResult{}, err
	}

	resp, err := http.Post("http://localhost:8000/query", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return QueryResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return QueryResult{}, fmt.Errorf("error: %s", string(body))
	}

	var result QueryResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return QueryResult{}, err
	}

	return result, nil
}

func getQueryNames() ([]string, error) {
	resp, err := http.Get("http://localhost:8000/query_names")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: status code %d", resp.StatusCode)
	}

	var names []string
	if err := json.NewDecoder(resp.Body).Decode(&names); err != nil {
		return nil, err
	}

	return names, nil
}

func main() {
	a := app.New()
	w := a.NewWindow("Pharmacy App")

	queryNames, err := getQueryNames()
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	buttons := make([]fyne.CanvasObject, len(queryNames))
	for i, name := range queryNames {
		queryID := i + 1
		buttons[i] = widget.NewButton(name, func() {
			showParameterForm(w, queryID)
		})
	}

	createOrderBtn := widget.NewButton("Create Order", func() {
		showCreateOrderForm(w)
	})

	viewOrdersBtn := widget.NewButton("View Orders", func() {
		showOrders(w)
	})

	editOrderBtn := widget.NewButton("Edit Order", func() {
		showEditOrderForm(w)
	})

	deleteOrderBtn := widget.NewButton("Delete Order", func() {
		showDeleteOrderForm(w)
	})

	content := container.NewVBox(buttons...)
	content.Add(widget.NewLabel("Order Management"))
	content.Add(createOrderBtn)
	content.Add(viewOrdersBtn)
	content.Add(editOrderBtn)
	content.Add(deleteOrderBtn)

	w.SetContent(content)
	w.Resize(fyne.NewSize(400, 600))
	w.CenterOnScreen()
	w.ShowAndRun()
}

func showParameterForm(parent fyne.Window, queryID int) {
	paramWindow := fyne.CurrentApp().NewWindow("Query parameters")
	paramEntries := make(map[string]*widget.Entry)

	switch queryID {
	case 1:
		queryString := strconv.Itoa(queryID)
		getQueryResultWithParams(parent, queryString, nil)
		getQueryResultWithParams(parent, queryString+"_count", nil)
		paramWindow.Close()
	case 2:
		paramEntries["Тип"] = widget.NewEntry()
		paramEntries["Тип"].SetPlaceHolder("Введите тип медикамента или оставьте поле пустым")
		formItems := make([]*widget.FormItem, 0, len(paramEntries))
		for label, entry := range paramEntries {
			formItems = append(formItems, widget.NewFormItem(label, entry))
		}

		form := widget.NewForm(formItems...)
		form.SubmitText = "Выполнить"
		form.OnSubmit = func() {
			params := make(map[string]string)
			for label, entry := range paramEntries {
				params[label] = entry.Text
			}
			if params["Тип"] == "" {
				queryString := strconv.Itoa(queryID)
				getQueryResultWithParams(parent, queryString, nil)
				getQueryResultWithParams(parent, queryString+"_count", nil)
			} else {
				getQueryResultWithParams(parent, strconv.Itoa(queryID)+"_type", params)
				getQueryResultWithParams(parent, strconv.Itoa(queryID)+"_type_count", params)
			}
			paramWindow.Close()
		}
		paramWindow.SetContent(form)
		paramWindow.Resize(fyne.NewSize(500, 200))
		paramWindow.Show()
	case 3:
		paramEntries["Тип"] = widget.NewEntry()
		paramEntries["Тип"].SetPlaceHolder("Введите тип медикамента или оставьте поле пустым")
		formItems := make([]*widget.FormItem, 0, len(paramEntries))
		for label, entry := range paramEntries {
			formItems = append(formItems, widget.NewFormItem(label, entry))
		}

		form := widget.NewForm(formItems...)
		form.SubmitText = "Выполнить"
		form.OnSubmit = func() {
			params := make(map[string]string)
			for label, entry := range paramEntries {
				params[label] = entry.Text
			}
			if params["Тип"] == "" {
				queryString := strconv.Itoa(queryID)
				getQueryResultWithParams(parent, queryString, nil)
			} else {
				getQueryResultWithParams(parent, strconv.Itoa(queryID)+"_type", params)
			}
			paramWindow.Close()
		}
		paramWindow.SetContent(form)
		paramWindow.Resize(fyne.NewSize(500, 200))
		paramWindow.Show()
	case 4:
		paramEntries["Вещество"] = widget.NewEntry()
		paramEntries["Вещество"].SetPlaceHolder("Введите название вещества")
		paramEntries["Начало периода"] = widget.NewEntry()
		paramEntries["Начало периода"].SetPlaceHolder("Укажите начало периода")
		paramEntries["Конец периода"] = widget.NewEntry()
		paramEntries["Конец периода"].SetPlaceHolder("Укажите конец периода")
		formItems := make([]*widget.FormItem, 0, len(paramEntries))
		for label, entry := range paramEntries {
			formItems = append(formItems, widget.NewFormItem(label, entry))
		}

		form := widget.NewForm(formItems...)
		form.SubmitText = "Выполнить"
		form.OnSubmit = func() {
			params := make(map[string]string)
			for label, entry := range paramEntries {
				params[label] = entry.Text
			}
			getQueryResultWithParams(parent, strconv.Itoa(queryID), params)
			paramWindow.Close()
		}
		paramWindow.SetContent(form)
		paramWindow.Resize(fyne.NewSize(500, 200))
		paramWindow.Show()
	case 6:
		queryString := strconv.Itoa(queryID)
		getQueryResultWithParams(parent, queryString, nil)
		paramWindow.Close()
	case 7:
		paramEntries["Тип"] = widget.NewEntry()
		paramEntries["Тип"].SetPlaceHolder("Введите категорию медикамента или оставьте поле пустым")
		formItems := make([]*widget.FormItem, 0, len(paramEntries))
		for label, entry := range paramEntries {
			formItems = append(formItems, widget.NewFormItem(label, entry))
		}

		form := widget.NewForm(formItems...)
		form.SubmitText = "Выполнить"
		form.OnSubmit = func() {
			params := make(map[string]string)
			for label, entry := range paramEntries {
				params[label] = entry.Text
			}
			if params["Тип"] == "" {
				queryString := strconv.Itoa(queryID)
				getQueryResultWithParams(parent, queryString, nil)
			} else {
				getQueryResultWithParams(parent, strconv.Itoa(queryID)+"_type", params)
			}
			paramWindow.Close()
		}
		paramWindow.SetContent(form)
		paramWindow.Resize(fyne.NewSize(500, 200))
		paramWindow.Show()
	case 8:
		queryString := strconv.Itoa(queryID)
		getQueryResultWithParams(parent, queryString, nil)
		getQueryResultWithParams(parent, queryString+"_count", nil)
		paramWindow.Close()
	case 9:
		queryString := strconv.Itoa(queryID)
		getQueryResultWithParams(parent, queryString, nil)
		paramWindow.Close()
	case 10:
		paramEntries["Тип"] = widget.NewEntry()
		paramEntries["Тип"].SetPlaceHolder("Введите категорию медикамента или оставьте поле пустым")
		formItems := make([]*widget.FormItem, 0, len(paramEntries))
		for label, entry := range paramEntries {
			formItems = append(formItems, widget.NewFormItem(label, entry))
		}

		form := widget.NewForm(formItems...)
		form.SubmitText = "Выполнить"
		form.OnSubmit = func() {
			params := make(map[string]string)
			for label, entry := range paramEntries {
				params[label] = entry.Text
			}
			if params["Тип"] == "" {
				queryString := strconv.Itoa(queryID)
				getQueryResultWithParams(parent, queryString, nil)
			} else {
				getQueryResultWithParams(parent, strconv.Itoa(queryID)+"_type", params)
			}
			paramWindow.Close()
		}
		paramWindow.SetContent(form)
		paramWindow.Resize(fyne.NewSize(500, 200))
		paramWindow.Show()
	case 12:
		paramEntries["Тип"] = widget.NewEntry()
		paramEntries["Тип"].SetPlaceHolder("Введите категорию медикамента")
		formItems := make([]*widget.FormItem, 0, len(paramEntries))
		for label, entry := range paramEntries {
			formItems = append(formItems, widget.NewFormItem(label, entry))
		}

		form := widget.NewForm(formItems...)
		form.SubmitText = "Выполнить"
		form.OnSubmit = func() {
			params := make(map[string]string)
			for label, entry := range paramEntries {
				params[label] = entry.Text
			}
			getQueryResultWithParams(parent, strconv.Itoa(queryID)+"_type", params)
			paramWindow.Close()
		}
		paramWindow.SetContent(form)
		paramWindow.Resize(fyne.NewSize(500, 200))
		paramWindow.Show()
	case 13:
		paramEntries["Тип"] = widget.NewEntry()
		paramEntries["Тип"].SetPlaceHolder("Введите название медикамента")
		formItems := make([]*widget.FormItem, 0, len(paramEntries))
		for label, entry := range paramEntries {
			formItems = append(formItems, widget.NewFormItem(label, entry))
		}

		form := widget.NewForm(formItems...)
		form.SubmitText = "Выполнить"
		form.OnSubmit = func() {
			params := make(map[string]string)
			for label, entry := range paramEntries {
				params[label] = entry.Text
			}
			fmt.Printf(params["Тип"])
			getQueryResultWithParams(parent, strconv.Itoa(queryID)+"_type", params)
			paramWindow.Close()
		}
		paramWindow.SetContent(form)
		paramWindow.Resize(fyne.NewSize(500, 200))
		paramWindow.Show()
	default:
		dialog.ShowError(fmt.Errorf("unknown query ID"), paramWindow)
		return
	}
}

func getQueryResultWithParams(parent fyne.Window, queryID string, params map[string]string) {
	queryParams := ""
	for key, value := range params {
		if queryParams != "" {
			queryParams += "&"
		}
		queryParams += fmt.Sprintf("%s=%s", key, value)
	}

	resp, err := http.Get(fmt.Sprintf("http://localhost:8000/queries/%s?%s", queryID, queryParams))
	if err != nil {
		dialog.ShowError(err, parent)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		dialog.ShowError(fmt.Errorf("error: %s", string(body)), parent)
		return
	}

	var result QueryResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		dialog.ShowError(err, parent)
		return
	}

	showResultTable(parent, result)
}

func showResultTable(parent fyne.Window, result QueryResult) {
	if len(result.Rows) == 0 {
		dialog.ShowInformation("Result", "No data found", parent)
		return
	}

	var data [][]string

	// Add column headers
	headers := make([]string, len(result.Columns))
	for i, col := range result.Columns {
		headers[i] = col
	}
	data = append(data, headers)

	// Add rows
	for _, row := range result.Rows {
		var rowData []string
		for _, val := range row {
			rowData = append(rowData, fmt.Sprintf("%v", val))
		}
		data = append(data, rowData)
	}

	resultTable := widget.NewTable(
		func() (int, int) { return len(data), len(data[0]) },
		func() fyne.CanvasObject {
			label := widget.NewLabel("------------------------------")
			label.Resize(fyne.NewSize(300, 50))
			return label
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			cell.(*widget.Label).SetText(data[id.Row][id.Col])
		},
	)

	resultWindow := fyne.CurrentApp().NewWindow("Query Result")
	resultWindow.SetContent(container.NewScroll(resultTable))
	resultWindow.Resize(fyne.NewSize(1400, 720))
	resultWindow.CenterOnScreen()
	resultWindow.Show()
}

func showCreateOrderForm(w fyne.Window) {
	customerIdEntry := widget.NewEntry()
	receiptIdEntry := widget.NewEntry()
	orderDateEntry := widget.NewEntry()
	productionDateEntry := widget.NewEntry()
	statusEntry := widget.NewEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Customer ID", Widget: customerIdEntry},
			{Text: "Receipt ID", Widget: receiptIdEntry},
			{Text: "Order Date (YYYY-MM-DD)", Widget: orderDateEntry},
			{Text: "Production Date (YYYY-MM-DD)", Widget: productionDateEntry},
			{Text: "Status", Widget: statusEntry},
		},
	}

	dialog.ShowForm("Create Order", "Create", "Cancel", form.Items, func(b bool) {
		if !b {
			return
		}
		customerID, err := strconv.Atoi(customerIdEntry.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid customer ID"), w)
			return
		}
		receiptID, err := strconv.Atoi(receiptIdEntry.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid receipt ID"), w)
			return
		}
		orderDate, err := time.Parse("2006-01-02", orderDateEntry.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid order date"), w)
			return
		}
		productionDate, err := time.Parse("2006-01-02", productionDateEntry.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid production date"), w)
			return
		}
		status := statusEntry.Text

		order := map[string]interface{}{
			"customer_id":     customerID,
			"receipt_id":      receiptID,
			"order_date":      orderDate.Format("2006-01-02"),
			"production_date": productionDate.Format("2006-01-02"),
			"status":          status,
		}

		data, err := json.Marshal(order)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		resp, err := http.Post("http://localhost:8000/orders", "application/json", bytes.NewBuffer(data))
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			body, _ := ioutil.ReadAll(resp.Body)
			dialog.ShowError(fmt.Errorf("error: %s", string(body)), w)
			return
		}

		dialog.ShowInformation("Success", "Order created successfully", w)
	}, w)
}

func showOrders(w fyne.Window) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8000/orders"))
	if err != nil {
		dialog.ShowError(err, w)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		dialog.ShowError(fmt.Errorf("error: %s", string(body)), w)
		return
	}

	var result QueryResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		dialog.ShowError(err, w)
		return
	}

	showResultTable(w, result)
}

func showEditOrderForm(w fyne.Window) {
	orderIdEntry := widget.NewEntry()
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Order ID", Widget: orderIdEntry},
		},
	}

	dialog.ShowForm("Edit Order", "Next", "Cancel", form.Items, func(b bool) {
		if !b {
			return
		}
		orderID, err := strconv.Atoi(orderIdEntry.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid order ID"), w)
			return
		}

		resp, err := http.Get(fmt.Sprintf("http://localhost:8000/orders/%d", orderID))
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := ioutil.ReadAll(resp.Body)
			dialog.ShowError(fmt.Errorf("error: %s", string(body)), w)
			return
		}

		var order Order
		if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
			dialog.ShowError(err, w)
			return
		}

		// Fetch customer, patient, doctor, and medicines data
		showEditCustomerForm(w, order.CustomerID)
		showEditPatientForm(w, order.ReceiptID)
		showEditDoctorForm(w, order.ReceiptID)
		showEditMedicineListForm(w, order.ReceiptID)
	}, w)
}

func showEditCustomerForm(w fyne.Window, customerID int) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8000/customers/%d", customerID))
	if err != nil {
		dialog.ShowError(err, w)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		dialog.ShowError(fmt.Errorf("error: %s", string(body)), w)
		return
	}

	var customer Customer
	if err := json.NewDecoder(resp.Body).Decode(&customer); err != nil {
		dialog.ShowError(err, w)
		return
	}

	surnameEntry := widget.NewEntry()
	surnameEntry.SetText(customer.Surname)
	nameEntry := widget.NewEntry()
	nameEntry.SetText(customer.Name)
	middleNameEntry := widget.NewEntry()
	middleNameEntry.SetText(customer.MiddleName)
	phoneNumberEntry := widget.NewEntry()
	phoneNumberEntry.SetText(customer.PhoneNumber)
	addressEntry := widget.NewEntry()
	addressEntry.SetText(customer.Address)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Surname", Widget: surnameEntry},
			{Text: "Name", Widget: nameEntry},
			{Text: "Middle Name", Widget: middleNameEntry},
			{Text: "Phone Number", Widget: phoneNumberEntry},
			{Text: "Address", Widget: addressEntry},
		},
	}

	dialog.ShowForm("Edit Customer", "Save", "Cancel", form.Items, func(b bool) {
		if !b {
			return
		}
		customer.Surname = surnameEntry.Text
		customer.Name = nameEntry.Text
		customer.MiddleName = middleNameEntry.Text
		customer.PhoneNumber = phoneNumberEntry.Text
		customer.Address = addressEntry.Text

		data, err := json.Marshal(customer)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:8000/customers/%d", customerID), bytes.NewBuffer(data))
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := ioutil.ReadAll(resp.Body)
			dialog.ShowError(fmt.Errorf("error: %s", string(body)), w)
			return
		}

		dialog.ShowInformation("Success", "Customer updated successfully", w)
	}, w)
}

func showEditPatientForm(w fyne.Window, receiptID int) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8000/receipts/%d/patient", receiptID))
	if err != nil {
		dialog.ShowError(err, w)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		dialog.ShowError(fmt.Errorf("error: %s", string(body)), w)
		return
	}

	var patient Patient
	if err := json.NewDecoder(resp.Body).Decode(&patient); err != nil {
		dialog.ShowError(err, w)
		return
	}

	surnameEntry := widget.NewEntry()
	surnameEntry.SetText(patient.Surname)
	nameEntry := widget.NewEntry()
	nameEntry.SetText(patient.Name)
	middleNameEntry := widget.NewEntry()
	middleNameEntry.SetText(patient.MiddleName)
	birthdayEntry := widget.NewEntry()
	birthdayEntry.SetText(patient.Birthday)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Surname", Widget: surnameEntry},
			{Text: "Name", Widget: nameEntry},
			{Text: "Middle Name", Widget: middleNameEntry},
			{Text: "Birthday", Widget: birthdayEntry},
		},
	}

	dialog.ShowForm("Edit Patient", "Save", "Cancel", form.Items, func(b bool) {
		if !b {
			return
		}
		patient.Surname = surnameEntry.Text
		patient.Name = nameEntry.Text
		patient.MiddleName = middleNameEntry.Text
		patient.Birthday = birthdayEntry.Text

		data, err := json.Marshal(patient)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:8000/patients/%d", patient.ID), bytes.NewBuffer(data))
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := ioutil.ReadAll(resp.Body)
			dialog.ShowError(fmt.Errorf("error: %s", string(body)), w)
			return
		}

		dialog.ShowInformation("Success", "Patient updated successfully", w)
	}, w)
}

func showEditDoctorForm(w fyne.Window, receiptID int) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8000/receipts/%d/doctor", receiptID))
	if err != nil {
		dialog.ShowError(err, w)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		dialog.ShowError(fmt.Errorf("error: %s", string(body)), w)
		return
	}

	var doctor Doctor
	if err := json.NewDecoder(resp.Body).Decode(&doctor); err != nil {
		dialog.ShowError(err, w)
		return
	}

	surnameEntry := widget.NewEntry()
	surnameEntry.SetText(doctor.Surname)
	nameEntry := widget.NewEntry()
	nameEntry.SetText(doctor.Name)
	middleNameEntry := widget.NewEntry()
	middleNameEntry.SetText(doctor.MiddleName)
	specialtyEntry := widget.NewEntry()
	specialtyEntry.SetText(doctor.Specialty)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Surname", Widget: surnameEntry},
			{Text: "Name", Widget: nameEntry},
			{Text: "Middle Name", Widget: middleNameEntry},
			{Text: "Specialty", Widget: specialtyEntry},
		},
	}

	dialog.ShowForm("Edit Doctor", "Save", "Cancel", form.Items, func(b bool) {
		if !b {
			return
		}
		doctor.Surname = surnameEntry.Text
		doctor.Name = nameEntry.Text
		doctor.MiddleName = middleNameEntry.Text
		doctor.Specialty = specialtyEntry.Text

		data, err := json.Marshal(doctor)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:8000/doctors/%d", doctor.ID), bytes.NewBuffer(data))
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := ioutil.ReadAll(resp.Body)
			dialog.ShowError(fmt.Errorf("error: %s", string(body)), w)
			return
		}

		dialog.ShowInformation("Success", "Doctor updated successfully", w)
	}, w)
}

func showEditMedicineListForm(w fyne.Window, receiptID int) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8000/receipts/%d/medicines", receiptID))
	if err != nil {
		dialog.ShowError(err, w)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		dialog.ShowError(fmt.Errorf("error: %s", string(body)), w)
		return
	}

	var medicines []Medicine
	if err := json.NewDecoder(resp.Body).Decode(&medicines); err != nil {
		dialog.ShowError(err, w)
		return
	}

	medicineEntries := make([]*widget.Entry, len(medicines))
	formItems := make([]*widget.FormItem, len(medicines))

	for i, medicine := range medicines {
		entry := widget.NewEntry()
		entry.SetText(medicine.Name)
		medicineEntries[i] = entry
		formItems[i] = &widget.FormItem{
			Text:   fmt.Sprintf("Medicine %d", i+1),
			Widget: entry,
		}
	}

	form := &widget.Form{
		Items: formItems,
	}

	dialog.ShowForm("Edit Medicine List", "Save", "Cancel", form.Items, func(b bool) {
		if !b {
			return
		}
		for i, entry := range medicineEntries {
			medicines[i].Name = entry.Text
		}

		data, err := json.Marshal(medicines)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:8000/receipts/%d/medicines", receiptID), bytes.NewBuffer(data))
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := ioutil.ReadAll(resp.Body)
			dialog.ShowError(fmt.Errorf("error: %s", string(body)), w)
			return
		}

		dialog.ShowInformation("Success", "Medicine list updated successfully", w)
	}, w)
}

func showDeleteOrderForm(w fyne.Window) {
	orderIdEntry := widget.NewEntry()
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Order ID", Widget: orderIdEntry},
		},
	}

	dialog.ShowForm("Delete Order", "Delete", "Cancel", form.Items, func(b bool) {
		if !b {
			return
		}
		orderID, err := strconv.Atoi(orderIdEntry.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid order ID"), w)
			return
		}

		req, err := http.NewRequest("DELETE", fmt.Sprintf("http://localhost:8000/orders/%d", orderID), nil)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			body, _ := ioutil.ReadAll(resp.Body)
			dialog.ShowError(fmt.Errorf("error: %s", string(body)), w)
			return
		}

		dialog.ShowInformation("Success", "Order deleted successfully", w)
	}, w)
}
