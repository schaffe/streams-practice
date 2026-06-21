package streams.practice.exercise2;

import java.util.Comparator;
import java.util.List;
import java.util.Map;

public class FlightTracker {

    public record Flight(String airline, String flightNumber, int durationMinutes) {
        @Override
        public String toString() {
            return String.format("%s (%s) - %d min", airline, flightNumber, durationMinutes);
        }
    }

    /**
     * Parses a list of flight data strings and returns a TreeMap sorted by duration (descending)
     * then airline name (alphabetically).
     * <p>
     * Each string must be in the format "Airline:FlightNumber:DurationInMinutes"
     * (e.g., "AirCanada:AC850:480").
     * <p>
     * Flights shorter than 90 minutes are filtered out.
     * If two flights have the same flight number, the later one overwrites the earlier.
     *
     * @param flightData list of flight data strings
     * @return TreeMap keyed by flight number, values are Flight objects, sorted by duration desc then airline asc
     */
    public Map<String, Flight> getFlightsSortedByDuration(List<String> flightData) {
        var res = flightData.stream()
                .map(s -> {
                    var sp = s.split(":");

                    return new Flight(sp[0], sp[1], Integer.parseInt(sp[2]));
                })
                .toList();

        System.out.println(res);

        // TODO: implement
        return null;
    }
}
