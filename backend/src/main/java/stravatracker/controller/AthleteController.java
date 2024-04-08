package stravatracker.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import stravatracker.model.Athlete;
import stravatracker.service.AthleteService;

import java.util.List;

@RestController
@RequestMapping("/api/v1/athlete")
public class AthleteController {


    private final AthleteService athleteService;

    @Autowired
    public AthleteController(AthleteService athleteService) {
        this.athleteService = athleteService;
    }

    @GetMapping("/athletes")
    public List<Athlete> getAllAthletes() {
        return this.athleteService.getAthletes();
    }
}
