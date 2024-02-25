package stravatracker.model;

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import jakarta.persistence.*;
import lombok.EqualsAndHashCode;

import java.util.ArrayList;
import java.util.List;


@Entity(name = "athletes")
@Table(name = "athletes",
        uniqueConstraints = {
                @UniqueConstraint(name = "athlete_firstname_lastname_unique", columnNames = {"firstname", "lastname"})
        }
)
//@ToString(exclude = "activities")
@EqualsAndHashCode(exclude = "activities")
public class Athlete {

    @Id
    @SequenceGenerator(name = "athlete_sequence", sequenceName = "athlete_sequence", allocationSize = 1)
    @GeneratedValue(strategy = GenerationType.IDENTITY, generator = "athlete_sequence")
    @Column(name = "id", updatable = false, unique = true, nullable = false)
    @JsonIgnore
    private long id;
    @Column(name = "firstname", columnDefinition = "TEXT", nullable = false)
    private String firstName;
    @Column(name = "lastname", columnDefinition = "TEXT", nullable = false)
    private String lastName;


    @OneToMany(mappedBy = "athlete", fetch = FetchType.LAZY, cascade = CascadeType.ALL)
    @JsonIgnoreProperties("activities")
    @JsonIgnore
    private final List<Activity> activities = new ArrayList<>();


    public Athlete() {

    }

    public Athlete(String firstName, String lastName) {
        this.firstName = firstName;
        this.lastName = lastName;
    }

    public Athlete(long id, String firstName, String lastName) {
        this.id = id;
        this.firstName = firstName;
        this.lastName = lastName;
    }


    public long getId() {
        return id;
    }

    public void setId(long id) {
        this.id = id;
    }

    public String getFirstName() {
        return firstName;
    }

    public void setFirstName(String firstName) {
        this.firstName = firstName;
    }

    public String getLastName() {
        return lastName;
    }

    public void setLastName(String lastName) {
        this.lastName = lastName;
    }

    public List<Activity> getActivities() {
        return activities;
    }

    @Override
    public String toString() {
        return "Athlete{" +
                "firstName='" + firstName + '\'' +
                ", lastName='" + lastName + '\'' +
                '}';
    }
}

