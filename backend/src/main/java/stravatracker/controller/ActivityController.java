package stravatracker.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;
import stravatracker.model.Activity;
import stravatracker.model.Athlete;
import stravatracker.requestsmodels.ActivitySumUpRequest;
import stravatracker.requestsmodels.ActivitySumUpResponse;
import stravatracker.requestsmodels.DateRange;
import stravatracker.service.ActivityService;

import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/v1/activity")
public class ActivityController {

    private final ActivityService activityService;

    @Autowired
    public ActivityController(ActivityService activityService) {
        this.activityService = activityService;
    }

    @GetMapping(value = "/activities")
    public List<Activity> getActivities() {
        return this.activityService.getActivities();
    }

    @GetMapping(value = "/activities/betweenDates")
    public List<Activity> getActivitiesBetweenDates(@RequestBody DateRange dateRange) {
        return this.activityService.getActivities(dateRange);
    }

    @GetMapping(value = "/athletes/sumUp")
    public Map<Athlete, ActivitySumUpResponse> getActivitiesSumUpByAthletes(
            @RequestBody ActivitySumUpRequest activitySumUpRequest) {
        return this.activityService.getActivitiesSumUpByAthletes(activitySumUpRequest.getSportType(),
                activitySumUpRequest.getDateRange());
    }
}
