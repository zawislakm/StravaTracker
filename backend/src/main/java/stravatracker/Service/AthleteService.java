package stravatracker.Service;

import com.fasterxml.jackson.databind.JsonNode;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import stravatracker.model.Athlete;
import stravatracker.repository.AthleteRepository;

import java.util.List;

@Service
public class AthleteService {

    private final AthleteRepository athleteRepository;

    @Autowired
    public AthleteService(AthleteRepository athleteRepository) {
        this.athleteRepository = athleteRepository;
    }

    public List<Athlete> getAthletes() {
        return athleteRepository.findAll();
    }


    public void addNewAthlete(Athlete athlete) {
        athleteRepository.save(athlete);
    }

    @Transactional
    public Athlete getOrCreateAthlete(JsonNode jsonData) {
        String firstname = jsonData.path("firstname").asText();
        String lastname = jsonData.path("lastname").asText();

        Athlete existingAthlete = athleteRepository.findByFirstNameAndLastName(firstname, lastname);

        if (existingAthlete != null) {
            return existingAthlete;
        } else {
            Athlete newAthlete = new Athlete();
            newAthlete.setFirstName(firstname);
            newAthlete.setLastName(lastname);

            return athleteRepository.save(newAthlete);
        }
    }
}
