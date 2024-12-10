package infrastructure

import (
	"fmt"
	"github.com/ITNS-LAB/gtfs-gorm/gtfsdb/domain/model"
	"github.com/ITNS-LAB/gtfs-gorm/gtfsdb/domain/repository"
	"github.com/ITNS-LAB/gtfs-gorm/gtfsjp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
	"path/filepath"
)

type gtfsScheduleRepository struct {
	Db  *gorm.DB
	Dsn string
}

func (g *gtfsScheduleRepository) ConnectDatabase() error {
	db, err := gorm.Open(postgres.Open(g.Dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	g.Db = db
	return nil
}

func (g *gtfsScheduleRepository) DisConnectDatabase() error {
	sqlDB, err := g.Db.DB()
	if err != nil {
		return err
	}
	if err := sqlDB.Close(); err != nil {
		return err
	}
	return nil
}

func (g *gtfsScheduleRepository) Migrate(shapeEx, shapeDetail bool) error {
	gtfsJp := model.GtfsJp{}
	err := g.Db.AutoMigrate(gtfsJp.Agency, gtfsJp.AgencyJp, gtfsJp.Routes, gtfsJp.Stops, gtfsJp.Calendar,
		gtfsJp.CalendarDates, gtfsJp.Trips, gtfsJp.StopTimes, gtfsJp.Shapes, gtfsJp.Frequencies, gtfsJp.Transfers,
		gtfsJp.FeedInfo, gtfsJp.FareAttribute, gtfsJp.FareRules, gtfsJp.Translations, gtfsJp.OfficeJp, gtfsJp.PatternJp)
	if err != nil {
		return err
	}
	if shapeEx {
		if err := g.Db.AutoMigrate(gtfsJp.ShapesEx); err != nil {
			return err
		}
	}
	if shapeDetail {
		if err := g.Db.AutoMigrate(gtfsJp.ShapesDetail); err != nil {
			return err
		}
	}
	return nil
}

func createGtfsJp[T any](filePath string, parser func(string) ([]T, error), db *gorm.DB) error {
	// ファイルのパース
	data, err := parser(filePath)
	if err != nil {
		slog.Warn(fmt.Sprintf("%sの読み込みができませんでした。ファイルが存在しない可能性があります。", filepath.Base(filePath)), "error", err, "path", filePath)
		return nil
	}

	// データベースへ挿入
	if err := db.CreateInBatches(&data, 1000).Error; err != nil {
		return fmt.Errorf("データベースへの挿入に失敗しました。%s", err)
	}
	slog.Info(fmt.Sprintf("%s の挿入が完了しました。", filepath.Base(filePath)))
	return nil
}

func (g *gtfsScheduleRepository) Create(gtfsPath string) error {
	if err := createGtfsJp(filepath.Join(gtfsPath, "agency.txt"), gtfsjp.ParseAgency, g.Db); err != nil {
		return err
	}
	if err := createGtfsJp(filepath.Join(gtfsPath, "agency_jp.txt"), gtfsjp.ParseAgencyJp, g.Db); err != nil {
		return err
	}
	if err := createGtfsJp(filepath.Join(gtfsPath, "calendar.txt"), gtfsjp.ParseCalendar, g.Db); err != nil {
		return err
	}
	if err := createGtfsJp(filepath.Join(gtfsPath, "calendar_dates.txt"), gtfsjp.ParseCalendarDates, g.Db); err != nil {
		return err
	}
	if err := createGtfsJp(filepath.Join(gtfsPath, "routes.txt"), gtfsjp.ParseRoutes, g.Db); err != nil {
		return err
	}
	if err := createGtfsJp(filepath.Join(gtfsPath, "stops.txt"), gtfsjp.ParseStops, g.Db); err != nil {
		return err
	}
	if err := createGtfsJp(filepath.Join(gtfsPath, "shapes.txt"), gtfsjp.ParseShapes, g.Db); err != nil {
		return err
	}
	if err := createGtfsJp(filepath.Join(gtfsPath, "trips.txt"), gtfsjp.ParseTrips, g.Db); err != nil {
		return err
	}
	if err := createGtfsJp(filepath.Join(gtfsPath, "stop_times.txt"), gtfsjp.ParseStopTimes, g.Db); err != nil {
		return err
	}
	if err := createGtfsJp(filepath.Join(gtfsPath, "transfers.txt"), gtfsjp.ParseTransfers, g.Db); err != nil {
		return err
	}
	if err := createGtfsJp(filepath.Join(gtfsPath, "frequencies.txt"), gtfsjp.ParseFrequencies, g.Db); err != nil {
		return err
	}
	if err := createGtfsJp(filepath.Join(gtfsPath, "fare_attributes.txt"), gtfsjp.ParseFareAttributes, g.Db); err != nil {
		return err
	}
	if err := createGtfsJp(filepath.Join(gtfsPath, "fare_rules.txt"), gtfsjp.ParseFareRules, g.Db); err != nil {
		return err
	}
	if err := createGtfsJp(filepath.Join(gtfsPath, "feed_info.txt"), gtfsjp.ParseFeedInfo, g.Db); err != nil {
		return err
	}
	if err := createGtfsJp(filepath.Join(gtfsPath, "translations.txt"), gtfsjp.ParseTranslations, g.Db); err != nil {
		return err
	}
	if err := createGtfsJp(filepath.Join(gtfsPath, "office_jp.txt"), gtfsjp.ParseOfficeJp, g.Db); err != nil {
		return err
	}
	if err := createGtfsJp(filepath.Join(gtfsPath, "pattern_jp.txt"), gtfsjp.ParsePatternJp, g.Db); err != nil {
		return err
	}
	return nil
}

func (g *gtfsScheduleRepository) CreateSchema(schema string) error {
	if err := g.Db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema)).Error; err != nil {
		return err
	}
	return nil
}

func (g *gtfsScheduleRepository) SetSchema(schema string) error {
	if err := g.Db.Exec(fmt.Sprintf("SET search_path TO %s, public", schema)).Error; err != nil {
		return fmt.Errorf("failed to set search_path: %w", err)
	}
	return nil
}

func (g *gtfsScheduleRepository) FindShapeIds() (shapeIds []string, err error) {
	g.Db.Table("shapes").Select("shape_id").Distinct("shape_id").Order("shape_id asc").Find(&shapeIds)
	return shapeIds, nil
}

func (g *gtfsScheduleRepository) FindShapes(shapeId string) (shapes []gtfsjp.Shape, err error) {
	g.Db.Table("shapes").Where("shape_id = ?", shapeId).Order("shape_pt_sequence asc").Find(&shapes)
	return shapes, nil
}

func (g *gtfsScheduleRepository) UpdateShapes(shapes []gtfsjp.Shape) error {
	tx := g.Db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for _, shape := range shapes {
		if result := tx.Model(&gtfsjp.Shape{}).
			Where("shape_id = ? AND shape_pt_sequence = ?", shape.ShapeId, shape.ShapePtSequence).
			Updates(shape); result.Error != nil {
			tx.Rollback() // エラーが発生したらロールバック
			return result.Error
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (g *gtfsScheduleRepository) FindTripsByShapeId(shapeId string) (trips []gtfsjp.Trip, err error) {
	g.Db.Table("trips").Where("shape_id = ?", shapeId).Find(&trips)
	return trips, nil
}

func (g *gtfsScheduleRepository) UpdateTrips(trips []gtfsjp.Trip) error {
	tx := g.Db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for _, trip := range trips {
		if result := tx.Model(&gtfsjp.Trip{}).
			Where("trip_id = ?", trip.TripId).
			Updates(trip); result.Error != nil {
			tx.Rollback() // エラーが発生したらロールバック
			return result.Error
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (g *gtfsScheduleRepository) FindTripIds() (tripIds []string, err error) {
	g.Db.Table("trips").Select("trip_id").Distinct("trip_id").Order("trip_id asc").Find(&tripIds)
	return tripIds, nil
}

func (g *gtfsScheduleRepository) FindShapesWithTripsByTripId(tripId string) (shapesEx []gtfsjp.ShapeEx, err error) {
	g.Db.Table("shapes").
		Select("trips.trip_id, trips.shape_id, shapes.shape_pt_lat, shapes.shape_pt_lon, shapes.shape_pt_sequence, shapes.shape_dist_traveled, NULL AS stop_id, shapes.geom").
		Joins("join trips on trips.shape_id = shapes.shape_id").
		Where("trip_id = ?", tripId).
		Order("shapes.shape_pt_sequence").
		Scan(&shapesEx)
	return shapesEx, nil
}

func (g *gtfsScheduleRepository) FindStopTimesByTripId(tripId string) (stopTimeWithLocations []model.StopTimeWithLocation, err error) {
	g.Db.Table("stop_times AS st").
		Select("st.trip_id, st.stop_id, st.stop_sequence, s.stop_lat, s.stop_lon").
		Joins("JOIN stops AS s ON s.stop_id = st.stop_id").
		Where("trip_id = ?", tripId).
		Order("stop_sequence").
		Scan(&stopTimeWithLocations)
	return stopTimeWithLocations, nil
}

func (g *gtfsScheduleRepository) CreateShapesEx(se []gtfsjp.ShapeEx) error {
	if err := g.Db.CreateInBatches(&se, 1000).Error; err != nil {
		return fmt.Errorf("データベースへの挿入に失敗しました。%s", err)
	}
	return nil
}

func (g *gtfsScheduleRepository) FetchShapes() (shapes []gtfsjp.Shape, err error) {
	g.Db.Find(&shapes)
	return shapes, nil
}

func (g *gtfsScheduleRepository) UpdateShapesEx(shapesEx []gtfsjp.ShapeEx) error {
	tx := g.Db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for _, shapeEx := range shapesEx {
		if result := tx.Model(&gtfsjp.ShapeEx{}).
			Where("trip_id = ? AND shape_id = ? AND shape_pt_sequence = ?", shapeEx.TripId, shapeEx.ShapeId, shapeEx.ShapePtSequence).
			Updates(shapeEx); result.Error != nil {
			tx.Rollback() // エラーが発生したらロールバック
			return result.Error
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (g *gtfsScheduleRepository) FetchShapesWithTrips() (shapesEx []gtfsjp.ShapeEx, err error) {
	g.Db.Table("shapes").
		Select("trips.trip_id, trips.shape_id, shapes.shape_pt_lat, shapes.shape_pt_lon, shapes.shape_pt_sequence, shapes.shape_dist_traveled, NULL AS stop_id, shapes.geom").
		Joins("join trips on trips.shape_id = shapes.shape_id").
		Order("trips.trip_id").
		Order("shapes.shape_pt_sequence").
		Scan(&shapesEx)
	return shapesEx, nil
}

func (g *gtfsScheduleRepository) UpdateStopTimes(stopTimes []gtfsjp.StopTime) error {
	tx := g.Db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for _, stopTime := range stopTimes {
		if result := tx.Model(&gtfsjp.StopTime{}).
			Where("trip_id = ? AND stop_id = ? AND stop_sequence = ?", stopTime.TripId, stopTime.StopId, stopTime.StopSequence).
			Updates(stopTime); result.Error != nil {
			tx.Rollback() // エラーが発生したらロールバック
			return result.Error
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (g *gtfsScheduleRepository) CreateShapeDetail(shapeDetails []gtfsjp.ShapeDetail) error {
	if err := g.Db.CreateInBatches(&shapeDetails, 1000).Error; err != nil {
		return fmt.Errorf("データベースへの挿入に失敗しました。%s", err)
	}
	return nil
}

func NewGtfsStaticRepository(dsn string) repository.GtfsScheduleRepository {
	var db *gorm.DB
	return &gtfsScheduleRepository{Db: db, Dsn: dsn}
}
