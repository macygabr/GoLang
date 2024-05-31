select orderuid, orders.tracknumber, entry, delivery.name, delivery.phone, delivery.Zip, delivery.City, delivery.Address, delivery.Region, delivery.Email,
payment.RequestID, payment.Currency, payment.Provider, payment.Amount, payment.paymentdt, payment.Bank, payment.deliverycost, payment.GoodsTotal, payment.CustomFee,
items.ChrtID, items.TrackNumber, items.Price, items.Rid, items.Name, items.Sale, items.Size, items.TotalPrice, items.NmID, items.Brand, items.Status,
orders.Locale, orders.InternalSignature, orders.CustomerID, orders.DeliveryService, orders.Shardkey, orders.SmID, orders.DateCreated, orders.OofShard
from orders 
    JOIN delivery ON orders.delivery_id = delivery.id 
    JOIN payment ON orders.payment_id = payment.id 
    JOIN items ON orders.items_id = items.id