package store

import "gondos/jetgen/gondos/table"

type tables struct {
	users     *table.UsersTable
	lists     *table.ListsTable
	listItems *table.ListItemsTable
}

func newTables(schema string) tables {
	return tables{
		users:     table.Users.FromSchema(schema),
		lists:     table.Lists.FromSchema(schema),
		listItems: table.ListItems.FromSchema(schema),
	}
}
