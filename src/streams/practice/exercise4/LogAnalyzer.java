package streams.practice.exercise4;

import java.util.*;
import java.util.stream.Collectors;

/**
 * Analyzes log lines and groups unique IP addresses by username.
 * <p>
 * Each log line is expected in the format:
 * "TIMESTAMP | USERNAME | ACTION | IP_ADDRESS"
 * <p>
 * IP addresses are sorted by subnet (first 3 octets) then by the last octet.
 */
public class LogAnalyzer {

    public record LogEntry(String timestamp, String username, String action, String ipAddress) {}

    public static final Comparator<String> IP_COMPARATOR = null; // TODO: implement Comparator that sorts by subnet then last octet

    /**
     * Groups log lines by USERNAME and collects unique IP addresses into a
     * sorted Set using a custom Comparator that sorts by subnet then last octet.
     *
     * @param logLines list of log line strings
     * @return map of username to sorted set of unique IP addresses
     */
    public Map<String, Set<String>> getIpAddressesByUser(List<String> logLines) {
        // TODO: implement
        return null;
    }

    /**
     * Parses a single log line into a LogEntry.
     * Expected format: "TIMESTAMP | USERNAME | ACTION | IP_ADDRESS"
     *
     * @param line the log line to parse
     * @return a LogEntry populated from the line
     */
    public static LogEntry parseLogLine(String line) {
        // TODO: implement
        return null;
    }

    public static List<String> sampleLogs() {
        return List.of(
                "2024-01-15 10:30:00 | alice | LOGIN | 192.168.1.100",
                "2024-01-15 10:31:00 | alice | VIEW_PAGE | 192.168.1.100",
                "2024-01-15 10:32:00 | bob | LOGIN | 10.0.0.5",
                "2024-01-15 10:33:00 | bob | LOGOUT | 10.0.0.5",
                "2024-01-15 10:34:00 | alice | LOGOUT | 10.0.0.1",
                "2024-01-15 10:35:00 | bob | LOGIN | 192.168.1.200",
                "2024-01-15 10:36:00 | charlie | LOGIN | 172.16.0.10",
                "2024-01-15 10:37:00 | charlie | VIEW_PAGE | 172.16.0.20",
                "2024-01-15 10:38:00 | bob | VIEW_PAGE | 10.0.0.5",
                "2024-01-15 10:39:00 | bob | DOWNLOAD | 192.168.1.200"
        );
    }

    public static void main(String[] args) {
        // TODO: create log lines, call getIpAddressesByUser, print results
    }
}
