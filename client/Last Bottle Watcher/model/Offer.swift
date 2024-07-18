//
//  Offer.swift
//  Last Bottle Watcher
//
//  Created by Marquard, Dave on 7/15/24.
//

import Foundation

struct Offer: Codable, Identifiable {
    var id: String
    var name: String
    var varietal: String?
    var vintage: String?
    var country: String?
    var region: String?
    var appellation: String?
    var bottleSize: String?
    var price: Int
    var retail: Int?
    var bestWeb: Int?
    var image: String
    
    var imageURL : URL? {
        get {
            if image.starts(with: "//") {
                URL(string: "https:\(image)")!
            }
            else {
                URL(string: image)!
            }
        }
    }
}
