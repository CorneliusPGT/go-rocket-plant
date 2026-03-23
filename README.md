POST
/api/v1/orders
Создать новый заказ
Request body
{
  "user_uuid": "string",
  "items": [
    {
      "part_uuid": "string",
      "quantity": 0
    }
  ]
}

GET
/api/v1/orders/{order_uuid}

POST
/api/v1/orders/{order_uuid}/pay
Оплатить заказ
Request body
{
  "payment_method": "CARD"
}