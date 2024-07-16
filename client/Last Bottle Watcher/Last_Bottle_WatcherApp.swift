//
//  Last_Bottle_WatcherApp.swift
//  Last Bottle Watcher
//
//  Created by Marquard, Dave on 7/15/24.
//

import SwiftUI

@main
struct Last_Bottle_WatcherApp: App {
    @State private var offerModel = OfferModel()
    
    var body: some Scene {
        WindowGroup {
            ContentView()
                .environment(offerModel)
        }
    }
}
