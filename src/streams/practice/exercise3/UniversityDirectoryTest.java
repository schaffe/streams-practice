package streams.practice.exercise3;

import static org.junit.jupiter.api.Assertions.*;

import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

class UniversityDirectoryTest {

    private UniversityDirectory directory;
    private List<UniversityDirectory.Department> departments;

    @BeforeEach
    void setUp() {
        var alice = new UniversityDirectory.Professor("Alice", Set.of("Java", "Python", "JavaScript"));
        var bob = new UniversityDirectory.Professor("Bob", Set.of("Java", "C++"));
        var carol = new UniversityDirectory.Professor("Carol", Set.of("Python", "Haskell", "Rust"));
        var dave = new UniversityDirectory.Professor("Dave", Set.of("Java", "TypeScript", "Go"));

        var cs = new UniversityDirectory.Department("Computer Science", List.of(alice, bob));
        var se = new UniversityDirectory.Department("Software Engineering", List.of(carol, dave));

        departments = List.of(cs, se);
        directory = new UniversityDirectory();
    }

    @Test
    void testGetAllLanguagesAlphabetized() {
        List<String> result = directory.getAllLanguagesAlphabetized(departments);
        List<String> expected = List.of("C++", "Go", "Haskell", "Java", "JavaScript", "Python", "Rust", "TypeScript");
        assertEquals(expected, result);
    }

    @Test
    void testCountProfessorsByLanguage() {
        Map<String, Long> result = directory.countProfessorsByLanguage(departments);
        assertEquals(3L, result.get("Java"));
        assertEquals(2L, result.get("Python"));
        assertEquals(1L, result.get("C++"));
        assertEquals(1L, result.get("Go"));
        assertEquals(1L, result.get("Haskell"));
        assertEquals(1L, result.get("JavaScript"));
        assertEquals(1L, result.get("Rust"));
        assertEquals(1L, result.get("TypeScript"));
        assertEquals(8, result.size());
    }

    @Test
    void testGetLanguagesByPopularity() {
        Map<String, Long> result = directory.getLanguagesByPopularity(departments);

        List<Map.Entry<String, Long>> entries = List.copyOf(result.entrySet());

        assertEquals("Java", entries.get(0).getKey());
        assertEquals(3L, entries.get(0).getValue());

        assertEquals("Python", entries.get(1).getKey());
        assertEquals(2L, entries.get(1).getValue());

        for (int i = 2; i < entries.size(); i++) {
            assertEquals(1L, entries.get(i).getValue().longValue());
        }

        assertTrue(result instanceof LinkedHashMap);
    }
}
