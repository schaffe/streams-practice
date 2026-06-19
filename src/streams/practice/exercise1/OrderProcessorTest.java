package streams.practice.exercise1;

import static org.junit.jupiter.api.Assertions.*;

import org.junit.jupiter.api.Test;
import java.util.List;

class OrderProcessorTest {

    @Test
    void testGetItemsSortedByPrice() {
        var order1 = new OrderProcessor.Order(List.of(
            new OrderProcessor.OrderItem("Cheap keychain", 1.50),
            new OrderProcessor.OrderItem("Luxury watch", 250.00),
            new OrderProcessor.OrderItem("Wallet", 45.00)
        ));

        var order2 = new OrderProcessor.Order(List.of(
            new OrderProcessor.OrderItem("Belt", 3.99),
            new OrderProcessor.OrderItem("Sunglasses", 15.00),
            new OrderProcessor.OrderItem("Hat", 15.00)
        ));

        var order3 = new OrderProcessor.Order(List.of(
            new OrderProcessor.OrderItem("Sunglasses", 15.00),
            new OrderProcessor.OrderItem("Scarf", 4.99),
            new OrderProcessor.OrderItem("Shoes", 89.99)
        ));

        List<OrderProcessor.Order> orders = List.of(order1, order2, order3);
        OrderProcessor processor = new OrderProcessor();
        List<OrderProcessor.OrderItem> result = processor.getItemsSortedByPrice(orders);

        assertEquals(5, result.size(), "Should have 5 items costing $5.00 or more");

        assertEquals("Luxury watch", result.get(0).productName());
        assertEquals(250.00, result.get(0).price(), 0.001);

        assertEquals("Shoes", result.get(1).productName());
        assertEquals(89.99, result.get(1).price(), 0.001);

        assertEquals("Wallet", result.get(2).productName());
        assertEquals(45.00, result.get(2).price(), 0.001);

        assertEquals("Hat", result.get(3).productName());
        assertEquals(15.00, result.get(3).price(), 0.001);

        assertEquals("Sunglasses", result.get(4).productName());
        assertEquals(15.00, result.get(4).price(), 0.001);

        assertTrue(result.get(3).productName().compareTo(result.get(4).productName()) <= 0,
                "Items with same price should be sorted alphabetically");
    }

    @Test
    void testGetItemsSortedByPriceAllCheap() {
        var order = new OrderProcessor.Order(List.of(
            new OrderProcessor.OrderItem("Sticker", 1.00),
            new OrderProcessor.OrderItem("Pin", 2.50),
            new OrderProcessor.OrderItem("Patch", 4.99)
        ));

        OrderProcessor processor = new OrderProcessor();
        List<OrderProcessor.OrderItem> result = processor.getItemsSortedByPrice(List.of(order));

        assertTrue(result.isEmpty(), "All items under $5.00 should be filtered out");
    }

    @Test
    void testGetItemsSortedByPriceEmptyOrders() {
        OrderProcessor processor = new OrderProcessor();
        List<OrderProcessor.OrderItem> result = processor.getItemsSortedByPrice(List.of());

        assertTrue(result.isEmpty(), "Empty orders list should produce empty result");
    }
}
