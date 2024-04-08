package stravatracker.model;

import com.fasterxml.jackson.annotation.JsonIgnore;
import jakarta.persistence.*;
import lombok.*;

import java.util.ArrayList;
import java.util.List;

@NoArgsConstructor
@AllArgsConstructor
@Getter
@Setter
@Entity(name = "sport_types")
@Table(name = "sport_types", uniqueConstraints = {
        @UniqueConstraint(name = "sport_type_unique", columnNames = {"type", "sport_type"})
})

@ToString(exclude = "activities")
@EqualsAndHashCode(exclude = "activities")
public class SportType {

    @Id
    @SequenceGenerator(name = "sport_type_sequence", sequenceName = "sport_type_sequence", allocationSize = 1)
    @GeneratedValue(strategy = GenerationType.IDENTITY, generator = "sport_type_sequence")
    @Column(name = "id", updatable = false, unique = true, nullable = false)
    @JsonIgnore
    private Long id;

    @Column(name = "type", columnDefinition = "TEXT", nullable = false)
    private String type;

    @Column(name = "sport_type", columnDefinition = "TEXT", nullable = false)
    private String sportType;

    @OneToMany(mappedBy = "sportType", fetch = FetchType.LAZY, cascade = CascadeType.ALL)
    @JsonIgnore
    private final List<Activity> activities = new ArrayList<>();

    public SportType(String type, String sportType) {
        this.type = type;
        this.sportType = sportType;
    }

}
