package streams.practice.exercise3;

import java.util.*;
import java.util.stream.Collectors;

public class UniversityDirectory {

    public record Professor(String name, Set<String> programmingLanguages) {}
    public record Department(String name, List<Professor> professors) {}

    /**
     * Use flatMap to get a stream of all programming languages from all professors
     * across all departments. Return a distinct, alphabetically sorted List of
     * language names.
     */
    public List<String> getAllLanguagesAlphabetized(List<Department> departments) {
        return departments.stream()
                .flatMap(d -> d.professors.stream())
                .flatMap(p -> p.programmingLanguages.stream())
                .distinct()
                .sorted()
                .toList();
    }

    /**
     * Create a Map where key is the programming language and value is the count
     * of professors who teach it.
     */
    public Map<String, Long> countProfessorsByLanguage(List<Department> departments) {

        return departments.stream()
                .flatMap(d -> d.professors.stream())
                .flatMap(p -> p.programmingLanguages.stream())
                .collect(Collectors.groupingBy(p -> p, Collectors.counting()));
    }

    /**
     * Same as above but the map must be sorted so the most popular language
     * (highest count) appears first. If two languages have the same count, sort
     * alphabetically. Hint: use a LinkedHashMap after sorting the entries.
     */
    public Map<String, Long> getLanguagesByPopularity(List<Department> departments) {
        return departments.stream()
                .flatMap(d -> d.professors.stream())
                .flatMap(p -> p.programmingLanguages.stream())
                .collect(Collectors.groupingBy(p -> p, Collectors.counting()))
                .entrySet()
                .stream()
                .sorted(Map.Entry.<String, Long>comparingByValue().reversed().thenComparing(Map.Entry::getKey))
                .peek(System.out::println)
                .collect(Collectors.toMap(e -> e.getKey(), e -> e.getValue(), (a, b) -> a, LinkedHashMap::new));
    }

    public static void main(String[] args) {
        // TODO: create departments and professors, call all 3 methods, print results
    }
}
