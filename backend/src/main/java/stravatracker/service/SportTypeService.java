package stravatracker.service;

import com.fasterxml.jackson.databind.JsonNode;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import stravatracker.model.SportType;
import stravatracker.repository.SportTypeRepository;

import java.util.ArrayList;
import java.util.List;

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


    public SportType getSportType(Long id) {
        return sportTypeRepository.findById(id).orElse(null);
    }


    public List<String> getSportTypeTypes() {


        List<String> types = new ArrayList<>();

        for (SportType sportType : sportTypeRepository.findAll()) {
            types.add(sportType.getType());
        }

        return types;
    }


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

        if (sportType.getType() != null) {
            return sportTypeRepository.findByType(sportType.getType());
        }

        if (sportType.getSportType() != null) {
            return sportTypeRepository.findBySportType(sportType.getSportType());
        }

        return null;
    }

    public SportType getSportType(String type) {
        return sportTypeRepository.findByType(type);
    }

    public SportType getSportType(String type, String sportType) {
        return sportTypeRepository.findByTypeAndSportType(type, sportType);
    }

    @Transactional
    public SportType getOrCreateSportType(JsonNode jsonData) {
        String type = jsonData.path("type").asText();
        String sport_type = jsonData.path("sport_type").asText();

        SportType existingSportType = sportTypeRepository.findByTypeAndSportType(type, sport_type);

        if (existingSportType != null) {
            return existingSportType;
        } else {
            SportType newSportType = new SportType();
            newSportType.setType(type);
            newSportType.setSportType(sport_type);

            return sportTypeRepository.save(newSportType);
        }
    }

}
