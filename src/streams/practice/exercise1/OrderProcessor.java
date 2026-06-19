package streams.practice.exercise1;

import java.util.Comparator;
import java.util.List;
import java.util.stream.Collectors;

public class OrderProcessor {

    public record OrderItem(String productName, double price) {}
    public record Order(List<OrderItem> items) {}

    /**
     * Processes a list of orders and returns all items sorted by price descending.
     * <ul>
     *   <li>Use flatMap to get a stream of all OrderItem objects across all orders</li>
     *   <li>Filter out items costing less than $5.00</li>
     *   <li>Sort by price descending; if same price, sort alphabetically by product name</li>
     *   <li>Collect into a List</li>
     * </ul>
     *
     * @param orders the list of orders to process
     * @return a list of items costing $5.00 or more, sorted by price descending then by name
     */
    public List<OrderItem> getItemsSortedByPrice(List<Order> orders) {
        // TODO: implement
        return null;
    }
}
