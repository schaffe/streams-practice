package streams.practice.exercise5;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class BuggyStream {

    /**
     * Supposed to return transaction amounts grouped by currency,
     * sorted alphabetically by currency, with amounts over $100 only.
     *
     * Find the TWO bugs and explain why they are wrong.
     */
    public static Map<String, List<Double>> processTransactions(List<String> transactionData) {
        Map<String, List<Double>> result = new HashMap<>();
        transactionData.stream()
            .forEach(t -> {
                String[] parts = t.split(",");
                String currency = parts[0];
                double amount = Double.parseDouble(parts[1]);
                if (amount > 100) {
                    result.computeIfAbsent(currency, k -> new ArrayList<>()).add(amount);
                }
            });
        return result;
    }

    public static void main(String[] args) {
        List<String> sampleData = List.of(
            "USD,250.00",
            "EUR,150.00",
            "GBP,50.00",
            "USD,75.00",
            "EUR,300.00",
            "JPY,5000.00",
            "USD,30.00",
            "GBP,200.00"
        );

        Map<String, List<Double>> result = processTransactions(sampleData);
        result.forEach((currency, amounts) -> System.out.println(currency + " -> " + amounts));
    }
}
