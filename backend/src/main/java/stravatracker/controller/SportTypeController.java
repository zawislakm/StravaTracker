package stravatracker.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import stravatracker.service.SportTypeService;
import stravatracker.model.SportType;

import java.util.List;

@RestController
@RequestMapping("/api/v1/sportType")
public class SportTypeController {

    private final SportTypeService sportTypeService;

    @Autowired
    public SportTypeController(SportTypeService sportTypeService) {
        this.sportTypeService = sportTypeService;
    }

    @GetMapping("/sportTypes")
    public List<SportType> getAllSportTypes() {
        return this.sportTypeService.getSportTypes();
    }

    @GetMapping("/types")
    public List<String> getSportTypeTypes() {
        return this.sportTypeService.getSportTypeTypes();
    }
}

