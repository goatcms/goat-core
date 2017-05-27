package orm

import "testing"

func TestFindByID(t *testing.T) {
	scope, err := newTestScope()
	if err != nil {
		t.Error(err)
		return
	}
	createTable, err := NewCreateTable(scope.table, scope.driver)
	if err != nil {
		t.Error(err)
		return
	}
	if err = createTable(scope.tx); err != nil {
		t.Error(err)
		return
	}
	insert, err := NewInsertWithID(scope.table, scope.driver)
	if err != nil {
		t.Error(err)
		return
	}
	inEntity := &TestEntity{10, "title1", "content1", "path1"}
	if err = insert(scope.tx, inEntity); err != nil {
		t.Error(err)
		return
	}
	findByID, err := NewFindByID(scope.table, scope.driver)
	if err != nil {
		t.Error(err)
		return
	}
	row, err := findByID(scope.tx, 10)
	if err != nil {
		t.Error(err)
		return
	}
	outEntity := &TestEntity{}
	if err = row.StructScan(outEntity); err != nil {
		t.Error(err)
		return
	}
	if outEntity.Content != inEntity.Content {
		t.Errorf("Content must be the same %v == %v", outEntity.Content, inEntity.Content)
		return
	}
	if outEntity.Title != inEntity.Title {
		t.Errorf("Title must be the same %v == %v", outEntity.Title, inEntity.Title)
		return
	}
	if outEntity.ID == 0 || inEntity.ID == 0 {
		t.Errorf("ID must be created %v != 0 && %v != 0", outEntity.ID, inEntity.ID)
		return
	}
	if outEntity.ID != inEntity.ID {
		t.Errorf("Title must be the same %v == %v", outEntity.ID, inEntity.ID)
		return
	}
	if outEntity.Image != inEntity.Image {
		t.Errorf("Image must be the same %v == %v", outEntity.Image, inEntity.Image)
		return
	}

}
