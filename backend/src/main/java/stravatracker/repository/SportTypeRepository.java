package stravatracker.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import stravatracker.model.SportType;

import java.util.List;

@Repository
public interface SportTypeRepository extends JpaRepository<SportType, Long> {
    SportType findByTypeAndSportType(String type, String sportType);

    List<SportType> findByType(String type);

    SportType findBySportType(String sportType);

}
