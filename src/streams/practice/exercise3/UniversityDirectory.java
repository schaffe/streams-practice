package streams.practice.exercise3;

import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;

public class UniversityDirectory {

    public record Professor(String name, Set<String> programmingLanguages) {}
    public record Department(String name, List<Professor> professors) {}

    /**
     * Use flatMap to get a stream of all programming languages from all professors
     * across all departments. Return a distinct, alphabetically sorted List of
     * language names.
     */
    public List<String> getAllLanguagesAlphabetized(List<Department> departments) {
        // TODO: implement
        return null;
    }

    /**
     * Create a Map where key is the programming language and value is the count
     * of professors who teach it.
     */
    public Map<String, Long> countProfessorsByLanguage(List<Department> departments) {
        // TODO: implement
        return null;
    }

    /**
     * Same as above but the map must be sorted so the most popular language
     * (highest count) appears first. If two languages have the same count, sort
     * alphabetically. Hint: use a LinkedHashMap after sorting the entries.
     */
    public Map<String, Long> getLanguagesByPopularity(List<Department> departments) {
        // TODO: implement
        return null;
    }

    public static void main(String[] args) {
        // TODO: create departments and professors, call all 3 methods, print results
    }
}
