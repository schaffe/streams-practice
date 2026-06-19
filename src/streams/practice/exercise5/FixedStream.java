package streams.practice.exercise5;

import java.util.List;
import java.util.Map;
import java.util.TreeMap;
import java.util.stream.Collectors;

public class FixedStream {

    /**
     * Fixed version: uses Stream API properly with collectors and TreeMap for sorting.
     * Transaction format: "Currency,Amount" e.g. "USD,250.00"
     */
    public static Map<String, List<Double>> processTransactions(List<String> transactionData) {
        // TODO: implement using stream(), filter(), collect()
        // Must use Collectors.groupingBy with TreeMap to ensure alphabetical order
        // Must use Collectors.mapping / Collectors.toList for the values
        return null;
    }

    public static void main(String[] args) {
        // TODO: create sample transaction data, call processTransactions, print results
    }
}
