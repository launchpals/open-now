//
//  DrawerViewController.swift
//  open-now-ios
//
//  Created by Robert Lin on 2019-01-27.
//  Copyright © 2019 launchpals. All rights reserved.
//

import UIKit
import Pulley

class DrawerViewController: UIViewController, PulleyDrawerViewControllerDelegate {
    
    var state = 0
    var currentViews = [UIView]()

    override func viewDidLoad() {
        super.viewDidLoad()

        // Do any additional setup after loading the view.
        setupViews()
        NotificationCenter.default.addObserver(self, selector: #selector(selectedStop(sender:)), name: NSNotification.Name(rawValue: "done"), object: nil)
    }
    
    func setupViews() {
        switch state {
        case 0:
            setupInitialViews()
        case 1:
            setupTable()
        case 2:
            setupDestination()
        default:
            return
        }
    }
    
    func setupInitialViews() {
        let titleLabel = UILabel()
        titleLabel.textColor = .white
        titleLabel.text = "I am..."
        titleLabel.font = UIFont.preferredFont(forTextStyle: .headline)
        titleLabel.sizeToFit()
        view.backgroundColor = #colorLiteral(red: 0.0431372549, green: 0.0431372549, blue: 0.0431372549, alpha: 1)
        view.alpha = 0.88
        
        view.addSubview(titleLabel)
        titleLabel.translatesAutoresizingMaskIntoConstraints = false
        titleLabel.centerXAnchor.constraint(equalTo: view.centerXAnchor, constant: 0).isActive = true
        titleLabel.topAnchor.constraint(equalTo: view.topAnchor, constant: 20).isActive = true
        
        let buttonLeft = UIButton()
        buttonLeft.setTitle("on the bus", for: .normal)
        buttonLeft.setTitleColor(UIColor.white, for: .normal)
        buttonLeft.backgroundColor = #colorLiteral(red: 0, green: 0.4, blue: 1, alpha: 1)
        buttonLeft.layer.cornerRadius = 20
        buttonLeft.layer.masksToBounds = true
        buttonLeft.addTarget(self, action: #selector(switchState), for: .touchUpInside)
//        buttonLeft.layer.borderWidth = 1
//        buttonLeft.layer.borderColor = UIColor.black.cgColor
        let buttonRight = UIButton()
        buttonRight.setTitle("walking", for: .normal)
        buttonRight.setTitleColor(UIColor.white, for: .normal)
//        buttonLeft.backgroundColor = #colorLiteral(red: 0, green: 0.4, blue: 1, alpha: 1)
        buttonRight.layer.cornerRadius = 20
        buttonRight.layer.masksToBounds = true
        buttonRight.layer.borderWidth = 2
        buttonRight.layer.borderColor = #colorLiteral(red: 0, green: 0.4, blue: 1, alpha: 1)
        let stackView = UIStackView(arrangedSubviews: [buttonLeft, buttonRight])
        stackView.axis = .horizontal
        stackView.spacing = 20
        stackView.distribution = .fillEqually
        view.addSubview(stackView)
        stackView.translatesAutoresizingMaskIntoConstraints = false
        stackView.centerXAnchor.constraint(equalTo: view.centerXAnchor, constant: 0).isActive = true
        stackView.topAnchor.constraint(equalTo: titleLabel.bottomAnchor, constant: 20).isActive = true
        stackView.widthAnchor.constraint(equalToConstant: 300).isActive = true
        stackView.heightAnchor.constraint(equalToConstant: 50).isActive = true
        
        currentViews = [titleLabel, stackView]
    }
    
    func setupTable() {
        let titleLabel = UILabel()
        titleLabel.textColor = .white
        titleLabel.text = "Which bus are you on?"
        titleLabel.font = UIFont.preferredFont(forTextStyle: .headline)
        titleLabel.sizeToFit()
        view.backgroundColor = #colorLiteral(red: 0.0431372549, green: 0.0431372549, blue: 0.0431372549, alpha: 1)
        view.alpha = 0.88
        
        view.addSubview(titleLabel)
        titleLabel.translatesAutoresizingMaskIntoConstraints = false
        titleLabel.centerXAnchor.constraint(equalTo: view.centerXAnchor, constant: 0).isActive = true
        titleLabel.topAnchor.constraint(equalTo: view.topAnchor, constant: 20).isActive = true
        
        let stackView = UIStackView()
        stackView.axis = .horizontal
        view.addSubview(stackView)
        stackView.spacing = 20
        stackView.translatesAutoresizingMaskIntoConstraints = false
        let busInfo = ["25", "41", "49", "70"]
        for bus in busInfo {
            let button = UIButton()
            button.setTitle(bus, for: .normal)
            button.setTitleColor(UIColor.white, for: .normal)
            button.layer.cornerRadius = 20
            button.layer.masksToBounds = true
            button.layer.borderWidth = 2
            button.layer.borderColor = #colorLiteral(red: 0, green: 0.4, blue: 1, alpha: 1)
            stackView.addArrangedSubview(button)
            button.addTarget(self, action: #selector(plotBus(sender:)), for: .touchUpInside)
        }
        stackView.distribution = .fillEqually
        stackView.centerXAnchor.constraint(equalTo: view.centerXAnchor, constant: 0).isActive = true
        stackView.topAnchor.constraint(equalTo: titleLabel.bottomAnchor, constant: 20).isActive = true
        stackView.widthAnchor.constraint(equalToConstant: 300).isActive = true
        stackView.heightAnchor.constraint(equalToConstant: 50).isActive = true
        currentViews = [titleLabel, stackView]
    }
    
    @objc func switchState() {
        UIView.animate(withDuration: 0.5) {
            for view in self.currentViews {
                view.alpha = 0
            }
            self.state += 1
            self.setupViews()
        }
    }
    
    func setupDestination() {
        let titleLabel = UILabel()
        titleLabel.textColor = .white
        titleLabel.text = "Choose a destination"
        titleLabel.font = UIFont.preferredFont(forTextStyle: .headline)
        titleLabel.sizeToFit()
        view.backgroundColor = #colorLiteral(red: 0.0431372549, green: 0.0431372549, blue: 0.0431372549, alpha: 1)
        view.alpha = 0.88
        
        view.addSubview(titleLabel)
        titleLabel.translatesAutoresizingMaskIntoConstraints = false
        titleLabel.centerXAnchor.constraint(equalTo: view.centerXAnchor, constant: 0).isActive = true
        titleLabel.topAnchor.constraint(equalTo: view.topAnchor, constant: 50).isActive = true
        currentViews = [titleLabel]
    }
    
    @objc func plotBus(sender: UIButton) {
        NotificationCenter.default.post(name: NSNotification.Name(rawValue: "bus"), object: sender.titleLabel!.text!)
        switchState()
    }
    
    func supportedDrawerPositions() -> [PulleyPosition] {
        return [PulleyPosition.partiallyRevealed, PulleyPosition.open]
    }
    
    func partialRevealDrawerHeight(bottomSafeArea: CGFloat) -> CGFloat {
        return bottomSafeArea + 120
    }
    
    @objc func selectedStop(sender: Notification) {
        
    }
    

    /*
    // MARK: - Navigation

    // In a storyboard-based application, you will often want to do a little preparation before navigation
    override func prepare(for segue: UIStoryboardSegue, sender: Any?) {
        // Get the new view controller using segue.destination.
        // Pass the selected object to the new view controller.
    }
    */

}
