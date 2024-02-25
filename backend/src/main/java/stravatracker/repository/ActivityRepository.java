package stravatracker.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import stravatracker.model.Activity;
import stravatracker.model.Athlete;
import stravatracker.model.SportType;

import java.time.LocalDate;
import java.util.List;

@Repository
public interface ActivityRepository extends JpaRepository<Activity, Long> {

    List<Activity> findByActivityDateBetween(LocalDate activityDate, LocalDate activityDate2);

    List<Activity> findByAthlete(Athlete athlete);

    List<Activity> findBySportType(SportType sportType);

    List<Activity> findBySportTypeAndAthlete(SportType sportType, Athlete athlete);

    List<Activity> findByAthleteAndActivityDateBetween(Athlete athlete, LocalDate start, LocalDate end);

    List<Activity> findByAthleteAndSportTypeAndActivityDateBetween(Athlete athlete, SportType sportType, LocalDate startDateAsLocalDate, LocalDate endDateAsLocalDate);

    List<Activity> findByAthleteAndSportType(Athlete athlete, SportType sportType);

    List<Activity> findBySportTypeAndActivityDateBetween(SportType sportType, LocalDate startDateAsLocalDate, LocalDate endDateAsLocalDate);

}
