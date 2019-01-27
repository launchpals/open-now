//
//  ViewController.swift
//  open-now-ios
//
//  Created by Yichen Cao on 2019-01-26.
//  Copyright © 2019 launchpals. All rights reserved.
//

import UIKit
import MapKit

class ViewController: UIViewController {
    
    let locationManager = CLLocationManager()
    var didSetup = false
    var latestLocation: CLLocation?
    var latestHeading: CLHeading?
    @IBOutlet weak var mapView: MKMapView!
    
    override func viewDidLoad() {
        super.viewDidLoad()
        
        setupMapView()
        setupLocationUpdates()
        setupGestureRecognizer()
    }
}

typealias ViewControllerLocationManager = ViewController
extension ViewControllerLocationManager: CLLocationManagerDelegate, MKMapViewDelegate {
    
    func setupMapView() {
        mapView.isRotateEnabled = true
        mapView.showsUserLocation = true
        mapView.showsCompass = true
        mapView.showsBuildings = false
        mapView._setShowsNightMode(true)
        mapView.showsPointsOfInterest = false
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
        
        if (!didSetup) {
//            let mapRegion = MKCoordinateRegion(center: latestLocation.coordinate, span: MKCoordinateSpan(latitudeDelta: 0.2, longitudeDelta: 0.2));
//            mapView.setRegion(mapRegion, animated: false)
//            let mapCamera = MKMapCamera(lookingAtCenter: latestLocation.coordinate, fromDistance: 100, pitch: 0, heading: 0)
//            mapView.setCamera(mapCamera, animated: false)
            mapView.setUserTrackingMode(.followWithHeading, animated: false)
            
            didSetup = true
        }
        self.latestLocation = latestLocation
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
    
    func plotRouteAt(coordinate: CLLocationCoordinate2D) {
        guard let latestLocation = latestLocation else {
            return
        }
        let request = MKDirections.Request()
        request.source = MKMapItem(placemark: MKPlacemark(coordinate: coordinate, addressDictionary: nil))
        request.destination = MKMapItem(placemark: MKPlacemark(coordinate: latestLocation.coordinate, addressDictionary: nil))
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
        renderer.strokeColor = UIColor(red: 17.0/255.0, green: 147.0/255.0, blue: 255.0/255.0, alpha: 1)
        renderer.lineWidth = 2.0
        return renderer
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
        plotRouteAt(coordinate: locationCoordinate)
    }
    
}

