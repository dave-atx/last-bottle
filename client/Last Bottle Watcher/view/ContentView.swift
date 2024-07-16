//
//  ContentView.swift
//  Last Bottle Watcher
//
//  Created by Marquard, Dave on 7/15/24.
//

import SwiftUI

struct ContentView: View {
    var body: some View {
        OfferListView()
    }
}

#Preview {
    ContentView().environment(OfferModel())
}
