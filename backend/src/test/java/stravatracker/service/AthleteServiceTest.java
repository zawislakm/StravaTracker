package stravatracker.service;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import jakarta.transaction.Transactional;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import stravatracker.model.Athlete;
import stravatracker.repository.AthleteRepository;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;

@SpringBootTest
class AthleteServiceTest {

    @Autowired
    private AthleteService athleteService;

    @Autowired
    private AthleteRepository athleteRepository;


    Athlete athlete1 = new Athlete("Professor", "K");
    Athlete athlete2 = new Athlete("Berlin", "O");
    Athlete athlete3 = new Athlete("Tokio", "A");

    @BeforeEach
    public void setUp() {
        athlete1 = athleteRepository.save(athlete1);
        athlete2 = athleteRepository.save(athlete2);
        athlete3 = athleteRepository.save(athlete3);
    }

    @AfterEach
    @Transactional
    void cleanUp() {
        athleteRepository.deleteAll();
    }

    @Test
    public void addNewAthleteTest() {
        Athlete athlete = new Athlete("Nairobi", "M");
        athleteService.addNewAthlete(athlete);

        Athlete athleteOutput = athleteRepository.findByFirstNameAndLastName(athlete.getFirstName(), athlete.getLastName());
        assertEquals(athleteOutput, athlete);
    }

    @Test
    public void getAllAthletes() {
        List<Athlete> expectedOutput = new ArrayList<>(Arrays.asList(
                athlete1,
                athlete2,
                athlete3
        ));

        List<Athlete> athleteListOutput = athleteService.getAthletes();

        assertEquals(expectedOutput, athleteListOutput);
    }

    @Test
    void getOrCreateSportType() {
        ObjectMapper mapper = new ObjectMapper();
        JsonNode existingJsonDataAthlete = mapper.createObjectNode()
                .put("firstname", athlete1.getFirstName())
                .put("lastname", athlete1.getLastName());


        Athlete newAthlete = new Athlete("Moscow", "V");

        JsonNode newJsonDataSportType = mapper.createObjectNode()
                .put("firstname", newAthlete.getFirstName())
                .put("lastname", newAthlete.getLastName());


        Athlete existingAthleteOutput = athleteService.getOrCreateAthlete(existingJsonDataAthlete);
        Athlete newAthleteOutput = athleteService.getOrCreateAthlete(newJsonDataSportType);


        assertEquals(existingAthleteOutput, athlete1);


        assertEquals(newAthleteOutput.getFirstName(), newAthlete.getFirstName());
        assertEquals(newAthleteOutput.getLastName(), newAthlete.getLastName());
    }

}