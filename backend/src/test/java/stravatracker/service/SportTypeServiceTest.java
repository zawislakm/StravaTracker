package stravatracker.service;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import jakarta.transaction.Transactional;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import stravatracker.model.SportType;
import stravatracker.repository.SportTypeRepository;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.stream.Collectors;

import static org.junit.jupiter.api.Assertions.assertEquals;

@SpringBootTest
class SportTypeServiceTest {

    @Autowired
    private SportTypeService sportTypeService;

    @Autowired
    private SportTypeRepository sportTypeRepository;

    SportType sportType1 = new SportType("Bike", "Bike");
    SportType sportType2 = new SportType("Bike", "MountainBike");
    SportType sportType3 = new SportType("Walk", "Mountain Walk");

    @BeforeEach
    public void setUp() {
        sportType1 = sportTypeRepository.save(sportType1);
        sportType2 = sportTypeRepository.save(sportType2);
        sportType3 = sportTypeRepository.save(sportType3);
    }

    @AfterEach
    @Transactional
    void cleanUp() {
        sportTypeRepository.deleteAll();
    }

    @Test
    void getSportTypes() {
        List<String> expectedTypes = new ArrayList<>(Arrays.asList(
                "Bike",
                "Walk"
        ));

        List<String> outputTypes = sportTypeService.getSportTypeTypes();

        List<String> sortedExpectedTypes = expectedTypes.stream().sorted().collect(Collectors.toList());
        List<String> sortedOutputTypes = outputTypes.stream().sorted().collect(Collectors.toList());

        assertEquals(sortedExpectedTypes, sortedOutputTypes);
    }

    @Test
    void addNewSportType() {
        SportType sportType = new SportType("Run", "Run");
        sportTypeService.addNewSportType(sportType);

        SportType sportTypeOutput = sportTypeService.getSportType(sportType);
        assertEquals(sportType, sportTypeOutput);
    }


    @Test
    void getOrCreateSportType() {
        ObjectMapper mapper = new ObjectMapper();
        JsonNode existingJsonDataSportType = mapper.createObjectNode()
                .put("type", sportType1.getType())
                .put("sport_type", sportType1.getSportType());


        SportType newSportType = new SportType("Run", "Mountain Run");

        JsonNode newJsonDataSportType = mapper.createObjectNode()
                .put("type", newSportType.getType())
                .put("sport_type", newSportType.getSportType());


        SportType existingSportTypeOutput = sportTypeService.getOrCreateSportType(existingJsonDataSportType);
        SportType newSportTypeOutput = sportTypeService.getOrCreateSportType(newJsonDataSportType);


        assertEquals(existingSportTypeOutput, sportType1);


        assertEquals(newSportTypeOutput.getType(), newSportType.getType());
        assertEquals(newSportTypeOutput.getSportType(), newSportType.getSportType());
    }


}