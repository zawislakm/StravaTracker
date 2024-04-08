package stravatracker.service;

import com.fasterxml.jackson.databind.JsonNode;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import stravatracker.model.SportType;
import stravatracker.repository.SportTypeRepository;

import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Set;

@Service
public class SportTypeService {

    private final SportTypeRepository sportTypeRepository;

    @Autowired
    public SportTypeService(SportTypeRepository sportTypeRepository) {
        this.sportTypeRepository = sportTypeRepository;
    }

    public List<SportType> getSportTypes() {
        return sportTypeRepository.findAll();
    }

    public void addNewSportType(SportType sportType) {
        sportTypeRepository.save(sportType);
    }


    public List<String> getSportTypeTypes() {
        Set<String> types = new HashSet<>();

        for (SportType sportType : sportTypeRepository.findAll()) {
            types.add(sportType.getType());
        }

        return new ArrayList<>(types);
    }

    //TODO refactor this method
    public SportType getSportType(SportType sportType) {

        if (sportType == null) {
            return null;
        }

        if (sportType.getId() != null) {
            return sportTypeRepository.findById(sportType.getId()).orElse(null);
        }

        if (sportType.getType() != null && sportType.getSportType() != null) {
            return sportTypeRepository.findByTypeAndSportType(sportType.getType(), sportType.getSportType());
        }


        if (sportType.getSportType() != null) {
            return sportTypeRepository.findBySportType(sportType.getSportType());
        }

        return null;
    }


    @Transactional
    public SportType getOrCreateSportType(JsonNode jsonData) {
        String type = jsonData.path("type").asText();
        String sport_type = jsonData.path("sport_type").asText();

        SportType existingSportType = sportTypeRepository.findByTypeAndSportType(type, sport_type);

        return (existingSportType != null) ? existingSportType : sportTypeRepository.save(new SportType(type, sport_type));
    }

}
