package streams.practice.exercise5;

import static org.junit.jupiter.api.Assertions.*;

import org.junit.jupiter.api.Test;

import java.util.List;
import java.util.Map;
import java.util.TreeMap;

class StreamRefactorTest {

    @Test
    void testProcessTransactions_ordersByCurrency() {
        List<String> data = List.of(
            "ZZZ,200.00",
            "AAA,150.00",
            "MMM,300.00"
        );
        Map<String, List<Double>> result = FixedStream.processTransactions(data);

        assertInstanceOf(TreeMap.class, result);
        List<String> keys = List.copyOf(result.keySet());
        assertEquals("AAA", keys.get(0));
        assertEquals("MMM", keys.get(1));
        assertEquals("ZZZ", keys.get(2));
    }

    @Test
    void testProcessTransactions_filtersUnder100() {
        List<String> data = List.of(
            "USD,200.00",
            "USD,50.00",
            "EUR,150.00",
            "EUR,25.00",
            "GBP,100.00"
        );
        Map<String, List<Double>> result = FixedStream.processTransactions(data);

        assertEquals(2, result.size());
        assertTrue(result.containsKey("USD"));
        assertTrue(result.containsKey("EUR"));
        assertFalse(result.containsKey("GBP"));

        assertEquals(1, result.get("USD").size());
        assertEquals(200.00, result.get("USD").get(0), 0.001);

        assertEquals(1, result.get("EUR").size());
        assertEquals(150.00, result.get("EUR").get(0), 0.001);
    }

    @Test
    void testProcessTransactions_groupsCorrectly() {
        List<String> data = List.of(
            "USD,250.00",
            "EUR,150.00",
            "USD,300.00",
            "EUR,175.00",
            "JPY,5000.00"
        );
        Map<String, List<Double>> result = FixedStream.processTransactions(data);

        assertEquals(3, result.size());
        assertEquals(2, result.get("USD").size());
        assertTrue(result.get("USD").contains(250.00));
        assertTrue(result.get("USD").contains(300.00));
        assertEquals(2, result.get("EUR").size());
        assertTrue(result.get("EUR").contains(150.00));
        assertTrue(result.get("EUR").contains(175.00));
        assertEquals(1, result.get("JPY").size());
        assertTrue(result.get("JPY").contains(5000.00));
    }
}
