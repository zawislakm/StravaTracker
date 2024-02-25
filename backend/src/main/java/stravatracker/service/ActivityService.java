package stravatracker.service;

import com.fasterxml.jackson.databind.JsonNode;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import stravatracker.model.Activity;
import stravatracker.model.Athlete;
import stravatracker.model.SportType;
import stravatracker.repository.ActivityRepository;
import stravatracker.requestsmodels.ActivitySumUpResponse;
import stravatracker.requestsmodels.DateRange;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

@Service
public class ActivityService {

    private final ActivityRepository activityRepository;

    private final AthleteService athleteService;

    private final SportTypeService sportTypeService;

    @Autowired
    public ActivityService(ActivityRepository activityRepository, AthleteService athleteService, SportTypeService sportTypeService) {
        this.activityRepository = activityRepository;
        this.athleteService = athleteService;
        this.sportTypeService = sportTypeService;
    }

    public void addNewActivity(Activity activity) {
        activityRepository.save(activity);
    }

    public void addNewActivity(JsonNode jsonData, Athlete athlete, SportType sportType) {
        Activity newActivity = new Activity(jsonData, athlete, sportType);
        activityRepository.save(newActivity);
    }


    public List<Activity> getActivities() {
        return activityRepository.findAll();
    }

    public List<Activity> getActivities(Athlete athlete) {
        return activityRepository.findByAthlete(athlete);
    }

    public List<Activity> getActivities(SportType sportType) {
        return activityRepository.findBySportType(sportType);
    }

    public List<Activity> getActivities(DateRange dateRange) {
        return activityRepository.findByActivityDateBetween(dateRange.getStartDate(), dateRange.getEndDate());
    }

    public List<Activity> getActivities(Athlete athlete, SportType sportType) {
        return activityRepository.findByAthleteAndSportType(athlete, sportType);
    }

    public List<Activity> getActivities(SportType sportType, DateRange dateRange) {
        return activityRepository.findBySportTypeAndActivityDateBetween(sportType, dateRange.getStartDate(), dateRange.getEndDate());
    }

    public List<Activity> getActivities(Athlete athlete, DateRange dateRange) {
        return activityRepository.findByAthleteAndActivityDateBetween(athlete, dateRange.getStartDate(), dateRange.getEndDate());
    }


    // TODO refactor this method
    public List<Activity> getActivities(Athlete athlete, SportType sportType, DateRange dateRange) {
        if (athlete != null && sportType != null && dateRange != null) {
            return activityRepository.findByAthleteAndSportTypeAndActivityDateBetween(athlete, sportType, dateRange.getStartDate(), dateRange.getEndDate());
        } else if (athlete == null && sportType != null && dateRange != null) {
            return this.getActivities(sportType, dateRange);
        } else if (athlete != null && sportType == null && dateRange != null) {
            return this.getActivities(athlete, dateRange);
        } else if (athlete != null && sportType != null && dateRange == null) {
            return this.getActivities(athlete, sportType);
        } else if (athlete == null && sportType == null && dateRange != null) {
            return this.getActivities(dateRange);
        } else if (athlete == null && sportType != null && dateRange == null) {
            return this.getActivities(sportType);
        } else if (athlete != null && sportType == null && dateRange == null) {
            return this.getActivities(athlete);
        } else {
            return this.getActivities();
        }
    }


    public Map<Athlete, ActivitySumUpResponse> getActivitiesSumUpByAthletes(SportType sportType, DateRange dateRange) {
        List<Athlete> athletes = athleteService.getAthletes();

        SportType sportTypeDatabase = sportTypeService.getSportType(sportType);
        Map<Athlete, ActivitySumUpResponse> sumUp = new HashMap<>();

        for (Athlete athlete : athletes) {
            List<Activity> activities = this.getActivities(athlete, sportTypeDatabase, dateRange);

            ActivitySumUpResponse activitySumUp = new ActivitySumUpResponse(activities);
            sumUp.put(athlete, activitySumUp);
        }

        return sumUp;
    }



}
