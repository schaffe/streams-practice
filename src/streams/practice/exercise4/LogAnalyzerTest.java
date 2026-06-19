package streams.practice.exercise4;

import static org.junit.jupiter.api.Assertions.*;

import org.junit.jupiter.api.Test;

import java.util.*;

class LogAnalyzerTest {

    @Test
    void testGroupByUser() {
        LogAnalyzer analyzer = new LogAnalyzer();
        Map<String, Set<String>> result = analyzer.getIpAddressesByUser(LogAnalyzer.sampleLogs());

        assertEquals(3, result.size());
        assertTrue(result.containsKey("alice"));
        assertTrue(result.containsKey("bob"));
        assertTrue(result.containsKey("charlie"));
    }

    @Test
    void testIpSortingOrder() {
        LogAnalyzer analyzer = new LogAnalyzer();
        Map<String, Set<String>> result = analyzer.getIpAddressesByUser(LogAnalyzer.sampleLogs());

        Set<String> bobIps = result.get("bob");
        List<String> bobIpList = new ArrayList<>(bobIps);

        assertEquals("10.0.0.5", bobIpList.get(0));
        assertEquals("192.168.1.200", bobIpList.get(1));

        Set<String> aliceIps = result.get("alice");
        List<String> aliceIpList = new ArrayList<>(aliceIps);

        assertEquals("10.0.0.1", aliceIpList.get(0));
        assertEquals("192.168.1.100", aliceIpList.get(1));

        Set<String> charlieIps = result.get("charlie");
        List<String> charlieIpList = new ArrayList<>(charlieIps);

        assertEquals("172.16.0.10", charlieIpList.get(0));
        assertEquals("172.16.0.20", charlieIpList.get(1));
    }

    @Test
    void testUniqueIpsOnly() {
        LogAnalyzer analyzer = new LogAnalyzer();
        Map<String, Set<String>> result = analyzer.getIpAddressesByUser(LogAnalyzer.sampleLogs());

        Set<String> aliceIps = result.get("alice");
        assertEquals(2, aliceIps.size());

        Set<String> bobIps = result.get("bob");
        assertEquals(2, bobIps.size());
    }
}
