package streams.practice.exercise2;

import org.junit.jupiter.api.Test;

import java.util.List;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.*;

class FlightTrackerTest {

    @Test
    void testGetFlightsSortedByDuration() {
        FlightTracker tracker = new FlightTracker();
        List<String> data = List.of(
            "Delta:DL404:320",
            "AirCanada:AC850:480",
            "United:UA120:180",
            "AirFrance:AF66:540"
        );
        Map<String, FlightTracker.Flight> result = tracker.getFlightsSortedByDuration(data);

        assertNotNull(result);
        assertEquals(4, result.size());

        List<FlightTracker.Flight> flights = List.copyOf(result.values());
        assertEquals(540, flights.get(0).durationMinutes());
        assertEquals(480, flights.get(1).durationMinutes());
        assertEquals(320, flights.get(2).durationMinutes());
        assertEquals(180, flights.get(3).durationMinutes());
    }

    @Test
    void testFilterShorterThan90Minutes() {
        FlightTracker tracker = new FlightTracker();
        List<String> data = List.of(
            "Delta:DL101:85",
            "United:UA120:180",
            "United:UA100:180",
            "Southwest:SWA45:60",
            "AirCanada:AC850:480"
        );
        Map<String, FlightTracker.Flight> result = tracker.getFlightsSortedByDuration(data);

        assertNotNull(result);
        assertEquals(3, result.size());
        assertFalse(result.containsKey("DL101"));
        assertFalse(result.containsKey("SWA45"));
    }
}
