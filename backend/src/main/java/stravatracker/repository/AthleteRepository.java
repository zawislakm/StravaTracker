package stravatracker.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import stravatracker.model.Athlete;

@Repository
public interface AthleteRepository extends JpaRepository<Athlete, Long> {
    Athlete findByFirstNameAndLastName(String firstname, String lastname);
}
