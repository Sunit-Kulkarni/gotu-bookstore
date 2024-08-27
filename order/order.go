// order/order.go
package order

import (
	"context"
	"encore.app/db"
	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
	"strconv"
)

//encore:api auth method=POST path=/orders
func CreateOrder(ctx context.Context, params *CreateOrderParams) (*CreateOrderResponse, error) {
	tx, err := db.Bookstoredb.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var orderID int
	err = tx.QueryRow(ctx, `
        INSERT INTO orders (user_id, created_at)
        VALUES ($1, NOW())
        RETURNING id
    `, strconv.Itoa(params.UserID)).Scan(&orderID)
	if err != nil {
		return nil, err
	}

	for _, item := range params.Items {
		_, err = tx.Exec(ctx, `
            INSERT INTO order_items (order_id, book_id, quantity)
            VALUES ($1, $2, $3)
        `, orderID, item.BookID, item.Quantity)
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &CreateOrderResponse{OrderID: orderID}, nil
}

//encore:api auth method=GET path=/orders
func GetOrderHistory(ctx context.Context) (*GetOrderHistoryResponse, error) {
	// Retrieve user id from JWT
	userID, _ := auth.UserID()
	authUserID, err := strconv.Atoi(string(userID))

	// First, fetch all orders for the user
	orderRows, err := db.Bookstoredb.Query(ctx, `
        SELECT id, created_at
        FROM orders
        WHERE user_id = $1
        ORDER BY created_at DESC
    `, userID)
	if err != nil {
		return nil, errs.WrapCode(err, errs.Internal, "failed to fetch orders")
	}
	defer orderRows.Close()

	var orders []OrderWithItems

	for orderRows.Next() {
		var order Order

		if err := orderRows.Scan(&order.ID, &order.CreatedAt); err != nil {
			return nil, errs.WrapCode(err, errs.Internal, "failed to scan order row")
		}

		order.UserID = authUserID
		orders = append(orders, OrderWithItems{Order: order, Items: []OrderItem{}})
	}

	if err := orderRows.Err(); err != nil {
		return nil, errs.WrapCode(err, errs.Internal, "error iterating over order rows")
	}

	// Now, fetch all order items for these orders
	itemRows, err := db.Bookstoredb.Query(ctx, `
        SELECT id, order_id, book_id, quantity
        FROM order_items
        WHERE order_id IN (SELECT id FROM orders WHERE user_id = $1)
        ORDER BY order_id, id
    `, userID)
	if err != nil {
		return nil, errs.WrapCode(err, errs.Internal, "failed to fetch order items")
	}
	defer itemRows.Close()

	itemMap := make(map[int][]OrderItem)
	for itemRows.Next() {
		var item OrderItem
		err := itemRows.Scan(&item.ID, &item.OrderID, &item.BookID, &item.Quantity)
		if err != nil {
			return nil, errs.WrapCode(err, errs.Internal, "failed to scan order item row")
		}
		itemMap[item.OrderID] = append(itemMap[item.OrderID], item)
	}

	if err := itemRows.Err(); err != nil {
		return nil, errs.WrapCode(err, errs.Internal, "error iterating over order item rows")
	}

	// Assign items to their respective orders
	for i, order := range orders {
		orders[i].Items = itemMap[order.Order.ID]
	}

	return &GetOrderHistoryResponse{Orders: orders}, nil
}
