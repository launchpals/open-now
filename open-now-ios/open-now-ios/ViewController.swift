//
//  ViewController.swift
//  open-now-ios
//
//  Created by Yichen Cao on 2019-01-26.
//  Copyright © 2019 launchpals. All rights reserved.
//

import UIKit
import MapKit

let address = "40.121.17.42:8081" // "34.73.17.110:8081" for old

class ViewController: UIViewController {
    
    let locationManager = CLLocationManager()
    var didSetup = false
    var latestLocation: CLLocation?
    var latestHeading: CLHeading?
    var latestPoiCenters = [CLLocationCoordinate2D]()
    weak var timer: Timer?
    lazy var client = OpenNow_CoreServiceClient.init(address: address, secure: false, arguments: [])
    @IBOutlet weak var mapView: MKMapView!
    var allStops = [String: [CLLocationCoordinate2D]]()
    
    override func viewDidLoad() {
        super.viewDidLoad()

        setupMapView()
        setupLocationUpdates()
//        setupGestureRecognizer()
        setupNotificationCenter()
    }
}

typealias ViewControllerLocationManager = ViewController
extension ViewControllerLocationManager: CLLocationManagerDelegate, MKMapViewDelegate {
    
    func setupMapView() {
        // disable interactions
        mapView.isRotateEnabled = false
        mapView.isScrollEnabled = false
        mapView.isPitchEnabled = false

        // set things to display
        mapView.showsUserLocation = true
        mapView.showsCompass = true
        mapView.showsBuildings = false
        mapView.showsPointsOfInterest = false
        mapView._setShowsNightMode(true)
        mapView.delegate = self
    }

    func setupLocationUpdates() {
        locationManager.delegate = self
        locationManager.desiredAccuracy = kCLLocationAccuracyBest
        locationManager.requestWhenInUseAuthorization()
        locationManager.startUpdatingLocation()
        locationManager.startUpdatingHeading()
    }
    
    func locationManager(_ manager: CLLocationManager, didUpdateLocations locations: [CLLocation]) {
        guard let latestLocation = locations.last else {
            return
        }
        
        self.latestLocation = latestLocation
        if (!didSetup) {
            latestPoiCenters.append(CLLocationCoordinate2D(latitude:  49.282963, longitude: -123.112358))
//            let mapCamera = MKMapCamera(lookingAtCenter: latestLocation.coordinate, fromDistance: 100, pitch: 70, heading: 0)
//            mapView.setCamera(mapCamera, animated: false)
            mapView.setUserTrackingMode(.followWithHeading, animated: false)
            fetchPOI()
            fetchTransit()
            didSetup = true
        }
//        updateMapTranslation()
    }

    /*
    func locationManager(_ manager: CLLocationManager, didUpdateHeading newHeading: CLHeading) {
        latestHeading = newHeading
        updateMapTranslation()
    }
    
    func updateMapTranslation() {
        guard let latestLocation = latestLocation, let latestHeading = latestHeading else {
            return
        }
        let mapCamera = MKMapCamera(lookingAtCenter: latestLocation.coordinate, fromDistance: 1000, pitch: 0, heading: latestHeading.magneticHeading)
        mapView.setCamera(mapCamera, animated: true)
    }*/
    
    func plotPOI(_ poi: OpenNow_Interest, _ source: CLLocationCoordinate2D?) {
        let coordinate = CLLocationCoordinate2D(latitude: poi.coordinates.latitude, longitude: poi.coordinates.longitude)
        plotRouteAt(target: coordinate, source: source)
        let poiAnnotation = MKPointAnnotation()
        poiAnnotation.title = poi.name
        poiAnnotation.coordinate = coordinate
        mapView.addAnnotation(poiAnnotation)
    }
    
    func distanceTo(_ c1: CLLocationCoordinate2D, _ c2: CLLocationCoordinate2D) -> CLLocationDistance {
        let p1 = MKMapPoint.init(c1)
        let p2 = MKMapPoint.init(c2)
        return p1.distance(to: p2)
    }
    
    func distanceTo(_ s: OpenNow_TransitStop, _ c: CLLocationCoordinate2D) -> CLLocationDistance {
        let coordinate = CLLocationCoordinate2D(latitude: s.coordinates.latitude, longitude: s.coordinates.longitude)
        let p1 = MKMapPoint.init(coordinate)
        let p2 = MKMapPoint.init(c)
        return p1.distance(to: p2)
    }
    
    func processTransit(_ stops: [OpenNow_TransitStop]) {
        // todo: fetch potential bus routes and display modal
        var routes: [String: [OpenNow_TransitStop]] = [:]
        var minDistanceToRoute: [String: CLLocationDistance] = [:]
        // var heading = locationManager.heading;
        for s in stops {
            let coordinate = CLLocationCoordinate2D(latitude: s.coordinates.latitude, longitude: s.coordinates.longitude)
            for r in s.routes {
                if let v = minDistanceToRoute[r] {
                    let dist = distanceTo(coordinate, latestLocation!.coordinate)
                    if v > dist {
                        minDistanceToRoute[r] = dist
                    }
                } else {
                     minDistanceToRoute[r] = distanceTo(coordinate, latestLocation!.coordinate)
                }
                
                if var stopList = routes[r] {
                    routes[r]!.append(s)
                } else {
                    routes[r] = [s]
                }
            }
        }
        let sortedRoutes = minDistanceToRoute.sorted { $0.1 < $1.1 }
        
        // Go through and render the first 3 stops for each route
        for (route, stops) in routes {
            if (allStops[route] == nil) {
                allStops[route] = [CLLocationCoordinate2D]()
            }
            let sortedStops = stops.sorted { distanceTo($0, latestLocation!.coordinate) < distanceTo($1, latestLocation!.coordinate) }
            
            var i = 0
            
            for s in sortedStops {
                print(route, i, distanceTo(s, latestLocation!.coordinate))
                if (i > 2) {
                    break
                }

                // TODO: we don't really need annotations for this
                let coordinate = CLLocationCoordinate2D(latitude: s.coordinates.latitude, longitude: s.coordinates.longitude)
                let stopAnnotation = MKPointAnnotation()
                stopAnnotation.title = s.routes.joined(separator: ", ")
                stopAnnotation.coordinate = coordinate
                mapView.addAnnotation(stopAnnotation)
                allStops[route]!.append(coordinate)
                
                i += 1
            }

        }
    }

    func mapView(_ mapView: MKMapView, viewFor annotation: MKAnnotation) -> MKAnnotationView? {
        guard annotation is MKPointAnnotation else { return nil }
        
        let identifier = "Annotation"
        var annotationView = mapView.dequeueReusableAnnotationView(withIdentifier: identifier)
        
        if annotationView == nil {
            annotationView = MKMarkerAnnotationView(annotation: annotation, reuseIdentifier: identifier)
            
            annotationView!.canShowCallout = true
            // TODO: HERE!
        } else {
            annotationView!.annotation = annotation
        }
        
        return annotationView
    }
    
    func plotRouteAt(target: CLLocationCoordinate2D, source: CLLocationCoordinate2D?) {
        var s = source
        if s == nil {
            guard let latestLocation = latestLocation else {
                return
            }
            s = latestLocation.coordinate
        }
        let request = MKDirections.Request()
        request.destination = MKMapItem(placemark: MKPlacemark(coordinate: target, addressDictionary: nil))
        request.source = MKMapItem(placemark: MKPlacemark(coordinate: s!, addressDictionary: nil))
        request.requestsAlternateRoutes = false
        request.transportType = .walking
        
        let directions = MKDirections(request: request)
        directions.calculate { [unowned self] response, error in
            guard let response = response, let route = response.routes.first, response.routes.count > 0 else { return }
            self.mapView.addOverlay(route.polyline)
//            self.mapView.setVisibleMapRect(route.polyline.boundingMapRect, animated: true)
        }
    }

    func mapView(_ mapView: MKMapView, rendererFor overlay: MKOverlay) -> MKOverlayRenderer {
        let renderer = MKPolylineRenderer(overlay: overlay)
        renderer.strokeColor = #colorLiteral(red: 0.2588235438, green: 0.7568627596, blue: 0.9686274529, alpha: 1)
        renderer.lineWidth = 4.0
        return renderer
    }
    
    func mapView(_ mapView: MKMapView, didSelect view: MKAnnotationView) {
        fetchPOI(coordinate: view.annotation!.coordinate)
    }
}

typealias ViewControllerGestures = ViewController
extension ViewControllerGestures: UIGestureRecognizerDelegate {
    
    func setupGestureRecognizer() {
        let tapgr = UITapGestureRecognizer(target: self, action: #selector(didTap(sender:)))
        mapView.addGestureRecognizer(tapgr)
    }
    
    @objc func didTap(sender: UITapGestureRecognizer) {
        let touchLocation = sender.location(in: mapView)
        let locationCoordinate = mapView.convert(touchLocation, toCoordinateFrom: mapView)
        plotRouteAt(target: locationCoordinate, source: nil)
    }
    
}

typealias ViewControllerFetch = ViewController
extension ViewControllerFetch {
    func fetchPOI() {
        guard let latestLocation = latestLocation else { return }
        fetchPOI(coordinate: latestLocation.coordinate)
    }
    
    func fetchPOI(coordinate: CLLocationCoordinate2D) {
        let coordinates = OpenNow_Coordinates.with {
            $0.latitude = coordinate.latitude
            $0.longitude = coordinate.longitude
        }
        let position = OpenNow_Position.with {
            $0.coordinates = coordinates
        }
        _ = try? client.getPointsOfInterest(position) { (pois, result) in
            guard let pois = pois?.interests else {
                return
            }
            for poi in pois {
                self.plotPOI(poi, coordinate)
            }
            NotificationCenter.default.post(name: NSNotification.Name(rawValue: "done"), object: pois.count)
        }
        
        if latestPoiCenters.count == 0 { return }
        for p in latestPoiCenters {
            let coordinates = OpenNow_Coordinates.with {
                $0.latitude = p.latitude
                $0.longitude = p.longitude
            }
            let position = OpenNow_Position.with {
                $0.coordinates = coordinates
            }
            _ = try? client.getPointsOfInterest(position) { (pois, result) in
                guard let pois = pois?.interests else {
                    return
                }
                for poi in pois {
                    self.plotPOI(poi, p)
                }
            }
        }
    }
    
    func fetchTransit() {
        guard let latestLocation = latestLocation else { return }
        let coordinates = OpenNow_Coordinates.with {
            $0.latitude = latestLocation.coordinate.latitude
            $0.longitude = latestLocation.coordinate.longitude
        }
        let position = OpenNow_Position.with {
            $0.coordinates = coordinates
        }
        _ = try? client.getTransitStops(position) { (stops, result) in
            guard let stops = stops?.stops else {
                return
            }
            self.processTransit(stops)
        }
    }
}

typealias ViewControllerNotificationCenter = ViewController
extension ViewControllerNotificationCenter {
    func setupNotificationCenter() {
        NotificationCenter.default.addObserver(self, selector: #selector(busSelected(bus:)), name: NSNotification.Name(rawValue: "bus"), object: nil)
    }
    
    @objc func busSelected(bus: NSNotification) {
        let busString = bus.object as! String
        print(busString + "selected!")
        mapView.removeOverlays(mapView.overlays)
        for annotation in mapView.annotations {
            if let title = annotation.title, let actualTitle = title, !actualTitle.components(separatedBy: ", ").contains("0" + busString) {
                mapView.removeAnnotation(annotation)
            }
        }
    }
}
