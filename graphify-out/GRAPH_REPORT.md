# Graph Report - app-listener  (2026-11-01)

## Corpus Check
- 565 files · ~1,800,552 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 7405 nodes · 33844 edges · 260 communities (218 shown, 42 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 3023 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `146a4824`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- go_pkg_strings
- testing.T
- IntegrationSuite
- op
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- IntegrationSuite
- StatInode
- mkdirs
- fakeGuardRepo
- Candidate
- policyEditor
- User
- buildOneGuard
- DaemonEvent
- IntegrationSuite
- guard_trust.bpf.c
- tabs.go
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- FileTarget
- guard.bpf.c
- ListUsers
- render.go
- gui/session.go
- fileedit_model_test.go
- highlight_test.go
- TestMonitorUseCaseLifecycle
- github.com/charmbracelet/bubbletea.Cmd
- recordFlow
- picker
- secret_probe/main.go
- IntegrationSuite
- Message
- wire.go
- EventType
- time.Time
- daemon.conf reference
- cursor.go
- monitorUnitTest
- .showSized
- filevault.go
- Guard
- stdio
- journal
- fadeScroll
- searchDirs
- protectedPage
- Client
- Network
- sync.Mutex
- .apply
- backend_test.go
- installPage
- engine
- Open
- TrustGuard
- serve.go
- SelectOrphans
- NetEventType
- unix_connect.c
- IntegrationSuite
- IntegrationSuite
- GuardEvent
- runNetworkGuard
- wipe
- LaunchRule
- planUpdaters
- fcntl
- filevault_test.go
- Vault
- codecheck/main.go
- logLine
- NetGuardEvent
- process_vm_readv.c
- os.File
- controlServer
- .filterEvent
- rcUser
- golang.org/x/sys/unix.Stat_t
- bpfstats/main.go
- exe_supersede.h
- loginServer
- RunWith
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- expandPlaceholders
- SanitizeText
- netinject.c
- README.md
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- time.Duration
- RefreshSection
- btrfsLayoutFrom
- Well-known bypass classes
- procstats.go
- jit_provenance.c
- ptrace_seize.c
- recordops.c
- runWizard
- Build & sign release assets (reusable workflow)
- BusClient
- github.com/charmbracelet/bubbles/textarea.Model
- hashFileShared
- UI
- guardUnitTest
- IntegrationSuite
- ConfinedEntry
- RunCmd
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- classify
- IntegrationSuite
- WebSocket /ws client
- newTestTarget
- inet_connect.c
- install.sh
- Config
- daemon: permanent protection with encryption at rest
- make test-integration (rootful Docker suite)
- vaultfs.go
- KernelDev
- .admitGated
- trace-app-libs.sh
- IntegrationSuite
- Monitor
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- TestCleanRel
- liveSession
- batchInodeRows
- Ledger
- serveWebSocket
- Daemon configuration
- runMonitor
- .existingWhitelist
- fakeStore
- runNetworkMonitor
- IntegrationSuite
- find.go
- runGenKey
- io/fs.FileInfo
- mc_tag.h
- Installation
- os.FileInfo
- graphify-refresh.sh
- .Layout
- update_test.go
- IntegrationSuite
- absPath
- feed.go
- shell
- runDaemon
- tempGrantSetup
- changelog.go
- SystemTrusted
- hashCache
- Limitations
- IntegrationSuite
- parseDaemonEvents
- net_tester/main.go
- plan_test.go
- sources
- IntegrationSuite
- fadescroll_test.go
- fileedit_symlink_test.go
- runView
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- Development
- github.com/spf13/cobra.Command
- Self-updating apps and the catalog
- run.go
- os.FileMode
- Existing
- edit-protected
- inspector_probe.c
- kernel-prof.sh
- .netIntegrityContainer
- pathCache
- bufio.Reader
- SecretStoreNotes
- compileGlobNames
- GuardInodeKey
- .populateWatchInodes
- monitor
- IntegrationSuite
- fetchReleases
- bpfLsmListed
- app-listener documentation
- .populateMaps
- WithBottomBar
- UnifiedDiff
- .awaitDaemonUpOrDead
- .TestRecord_GuardParity
- inodeChain
- runBackend
- changelog_test.go
- .eventTypeAllowed
- find
- ReadInstallCert
- listFacultatives
- editView
- supersede_probe.c
- sources_test.go
- bleedLayout
- mountTable
- deploy
- github.com/Virgula0/app-listener/tools/codecheck
- dialLiveSession
- network monitor: see what a program talks to
- inode_dev
- mutedText
- NewDumpHook
- Troubleshooting
- configedit_test.go
- GenerateConf
- run
- shQuote
- taint_victim.c
- NewEventFanout
- maintainer/plan.go
- testLimiter
- pinstate.go
- linkView
- writeWithin
- checkBtrfsKeys
- parseKernelRelease
- parseResize
- golang.org/x/sys/unix.Timespec
- VouchedDevs
- runRefreshCatalogOnly
- fileID
- copyXattrs
- MulticallError
- watchGroup

## God Nodes (most connected - your core abstractions)
1. `Config` - 147 edges
2. `Guard` - 125 edges
3. `GuardInodeKey` - 112 edges
4. `Load()` - 97 edges
5. `shell` - 97 edges
6. `absPath()` - 93 edges
7. `IntegrationSuite` - 83 edges
8. `fileEditModel` - 77 edges
9. `writeConfig()` - 75 edges
10. `IntegrationSuite` - 67 edges

## Surprising Connections (you probably didn't know these)
- `Modes` --references--> `monitor()`  [INFERRED]
  docs/README.md → integrationtests/secret_probe/main.go
- `Distributions` --references--> `monitor()`  [INFERRED]
  docs/compatibility.md → integrationtests/secret_probe/main.go
- `Files` --references--> `monitor()`  [INFERRED]
  docs/gui.md → integrationtests/secret_probe/main.go
- `Workspaces` --references--> `monitor()`  [INFERRED]
  docs/gui.md → integrationtests/secret_probe/main.go
- `eBPF at three levels` --references--> `monitor()`  [INFERRED]
  docs/how-it-works.md → integrationtests/secret_probe/main.go

## Import Cycles
- None detected.

## Communities (260 total, 42 thin omitted)

### Community 0 - "go_pkg_strings"
Cohesion: 0.09
Nodes (19): tempJournalRow, Progress, User, invalidPassword, recordTick, linkedFile, argvMatch, launchRule (+11 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (190): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestEveryInstallModeMapped(), TestForwardFlagSurface(), TestSubcommandsRegistered() (+182 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (7): configBinaryLines(), IntegrationSuite, probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig()

### Community 3 - "op"
Cohesion: 0.17
Nodes (13): clearExcept(), createTemp(), holds(), readlinkat(), removeDir(), ignoreNotExist(), openBeneath(), openDir() (+5 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (82): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+74 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.09
Nodes (51): prctl_gate(), guard_task_prctl(), admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader() (+43 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.11
Nodes (48): current_exe(), emit_dentry(), emit_event(), emit_event_kern(), emit_file_write(), emit_kern(), emit_unix_peer(), emit_user() (+40 more)

### Community 7 - "IntegrationSuite"
Cohesion: 0.24
Nodes (6): netMonitorEvent, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 8 - "StatInode"
Cohesion: 0.09
Nodes (14): ReplacementCheck, copySelf(), runTool(), WithEagerPopulate(), appendByte(), guardUnitTest, TestReadOnlyTreeInodeDoesNotFollowSymlinks(), SetReplacementCheck() (+6 more)

### Community 9 - "mkdirs"
Cohesion: 0.11
Nodes (43): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), newWatchPattern(), systemBinaryPatterns(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped() (+35 more)

### Community 10 - "fakeGuardRepo"
Cohesion: 0.06
Nodes (16): TempBinary, TemporaryGrant, TempRule, applyTempGrants(), closeTempBinaries(), TemporaryAccess, logTempGrant(), newExeTap() (+8 more)

### Community 11 - "Candidate"
Cohesion: 0.10
Nodes (30): nextConfig(), appendSectionsAndEdit(), collectDiffAdditions(), TestUncoveredCandidates(), uncoveredCandidates(), selectAndEditConfig(), addManualDirectories(), editConfig() (+22 more)

### Community 12 - "policyEditor"
Cohesion: 0.11
Nodes (16): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), eventsOrAll(), netWarnings() (+8 more)

### Community 13 - "User"
Cohesion: 0.11
Nodes (27): entryWriters(), reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), bunEntryConfigured(), bunLaunchers(), bunLaunchersFor() (+19 more)

### Community 14 - "buildOneGuard"
Cohesion: 0.06
Nodes (55): CheckBPFLSM(), CheckEBPF(), backingDeviceUnion(), buildGuards(), buildOneGuard(), eventFilterOptions(), newPinGeneration(), prepareDaemonStart() (+47 more)

### Community 15 - "DaemonEvent"
Cohesion: 0.12
Nodes (24): UIDResolver, NewUIDResolver(), drainEphemeral(), runHeadless(), writeEvent(), admitEvent(), foldable(), newGateLogLimiter() (+16 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (8): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (64): argv_env(), argv_matched(), argv_rule_at(), argv_word(), caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key() (+56 more)

### Community 18 - "tabs.go"
Cohesion: 0.03
Nodes (25): appTheme, binaryRow, clickPanel, iconTile, iconTileRenderer, navItem, navItemRenderer, panel (+17 more)

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (11): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, fileEditMode (+3 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.14
Nodes (51): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+43 more)

### Community 22 - "Resource"
Cohesion: 0.05
Nodes (13): fakeDaemon, Resource, GuardRepository, Vault, ConcurrencyLimit(), encryptionRootOf(), DaemonUseCase, TestRepeatLogOncePerErrorAndOnRecovery() (+5 more)

### Community 23 - "auth.go"
Cohesion: 0.12
Nodes (42): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+34 more)

### Community 24 - "FileTarget"
Cohesion: 0.11
Nodes (14): guardPill(), Files, FileSpec, FileTarget, Policy, guardRow(), newFileTarget(), NormalizeSpec() (+6 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.07
Nodes (95): exe_refused(), add_inode_to_guard(), admit_port(), btrfs_copy_gate(), bus_key(), check_and_emit_args(), chmod_only_drops_write(), code_suspect() (+87 more)

### Community 26 - "ListUsers"
Cohesion: 0.07
Nodes (33): releaseBuild(), Run(), grantNewBinaries(), cleanOrphanedFscrypt(), askDecrypt(), cleanOrphanedMetadata(), confirmKeyDeletion(), decryptStep() (+25 more)

### Community 27 - "render.go"
Cohesion: 0.12
Nodes (13): lineStyles, viewportState, baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces(), setGutter() (+5 more)

### Community 28 - "gui/session.go"
Cohesion: 0.07
Nodes (40): ancestorEnv(), invokingHome(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField(), TestDisplayValue() (+32 more)

### Community 29 - "fileedit_model_test.go"
Cohesion: 0.20
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 30 - "highlight_test.go"
Cohesion: 0.12
Nodes (26): New(), SetText(), TestBackspaceAndDeleteRemoveAnIndentStep(), TestCursorDoesNotBlink(), TestEditKeys(), TestEnterKeepsIndentation(), TestEnterPastNinetyNineLines(), TestSetTextStartsAtFirstLine() (+18 more)

### Community 31 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.18
Nodes (12): NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestMonitorUseCaseStartError() (+4 more)

### Community 32 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.04
Nodes (42): recordTickCmd(), recordModel, diffModel, editorModel, formatGuardEventLine(), sanitizeTerminalText(), sanitizeTerminalTexts(), formatDecision() (+34 more)

### Community 33 - "recordFlow"
Cohesion: 0.05
Nodes (57): encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), lines(), offerBackupDeletion(), pickBinaries() (+49 more)

### Community 34 - "picker"
Cohesion: 0.08
Nodes (30): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, openPicker(), pickerTree() (+22 more)

### Community 35 - "secret_probe/main.go"
Cohesion: 0.20
Nodes (27): busPath(), connexec(), connFD(), createItem(), dial(), dupkey(), execInto(), failf() (+19 more)

### Community 37 - "Message"
Cohesion: 0.10
Nodes (14): checkSecret(), checkTrustLine(), trustRows(), File, TrustItem, backend, backendOp, backendState (+6 more)

### Community 38 - "wire.go"
Cohesion: 0.12
Nodes (19): dictFrame(), TestClientAuth(), TestDuplicateDictKeyRefused(), TestReadFrame(), TestSASL(), clientAuth(), decodeFrame(), distinctDictKeys() (+11 more)

### Community 39 - "EventType"
Cohesion: 0.06
Nodes (20): selfEventsFor(), snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting() (+12 more)

### Community 40 - "time.Time"
Cohesion: 0.14
Nodes (10): LockedError, isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet, sweepID (+2 more)

### Community 41 - "daemon.conf reference"
Cohesion: 0.08
Nodes (24): A complete example, `allow_env`, `allow_lib`, Binary lines and events, daemon.conf reference, `deny_args`, `warn_args`, `deny_env`, `warn_env`, `[electron_apps]` (+16 more)

### Community 42 - "cursor.go"
Cohesion: 0.08
Nodes (20): buttonCursor, handCheckGroup, handSelect, handCursor(), handCursors(), handCursorsOn(), newButtonCursor(), newCheckGroup() (+12 more)

### Community 44 - ".showSized"
Cohesion: 0.10
Nodes (13): binaryList, installUI, multiAnswer, scrollable(), selectable(), newFadeScroll(), ask(), diffText() (+5 more)

### Community 45 - "filevault.go"
Cohesion: 0.21
Nodes (15): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), newFileVaultAEAD(), openFileVault(), openFileVaultWithMasterKey() (+7 more)

### Community 46 - "Guard"
Cohesion: 0.06
Nodes (33): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), binaryVerifyState, deferredBinary, exeRow, keptEntries (+25 more)

### Community 48 - "journal"
Cohesion: 0.14
Nodes (6): fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal

### Community 49 - "fadeScroll"
Cohesion: 0.10
Nodes (5): fadeBar, fadeBarRenderer, fadeScroll, fadeScrollRenderer, fadeScrollOf()

### Community 50 - "searchDirs"
Cohesion: 0.18
Nodes (10): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), ldConfInclude(), ldSoConfDirs(), LdSoPreloadPaths(), resolveSoname() (+2 more)

### Community 51 - "protectedPage"
Cohesion: 0.05
Nodes (29): backendLink, configView, forwardForm, forwardView, handRadioGroup, protectedPage, recordReview, recordView (+21 more)

### Community 52 - "Client"
Cohesion: 0.09
Nodes (14): Client, childEnv(), launch(), Start(), TestChildEnvIsAllowlisted(), TestHelperBackend(), TestReaderCapsLines(), Reader (+6 more)

### Community 53 - "Network"
Cohesion: 0.12
Nodes (8): netGuardPill(), NetworkGuardRepository, clonePolicy(), NetEngines, NetExtra, NetPolicy, Network, NetGuardConfig

### Community 54 - "sync.Mutex"
Cohesion: 0.11
Nodes (9): fyneLogFilter, Gate, Listen(), peerCred(), readFrame(), Bus, queued, senderQueue (+1 more)

### Community 55 - ".apply"
Cohesion: 0.11
Nodes (12): Prereq, Apply(), DeleteBackup(), Outcome, ProgressFunc, TestDeleteBackupOnlyItsOwn(), applier, fakeVault (+4 more)

### Community 56 - "backend_test.go"
Cohesion: 0.23
Nodes (22): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+14 more)

### Community 57 - "installPage"
Cohesion: 0.05
Nodes (33): anyDaemonRunning(), actionGate, binSource, handButton, handCheck, installAction, installedBinary, InstallMode (+25 more)

### Community 58 - "engine"
Cohesion: 0.04
Nodes (26): engine, resInfo, taintOwner, TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), pathWithin(), TestPathWithinMatchesPrefixRule() (+18 more)

### Community 59 - "Open"
Cohesion: 0.11
Nodes (26): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+18 more)

### Community 60 - "TrustGuard"
Cohesion: 0.05
Nodes (26): tempGrant, tempMaskOp, tempRow, trustHook, trustRows, droppableFreed(), freedSuperseded(), PruneSuperseded() (+18 more)

### Community 61 - "serve.go"
Cohesion: 0.17
Nodes (18): basicAuth(), isInteractiveTerminal(), newServeListener(), newServeServer(), relayServeSignals(), runServe(), securityHeaders(), Serve() (+10 more)

### Community 63 - "SelectOrphans"
Cohesion: 0.23
Nodes (16): appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors(), SelectOrphans() (+8 more)

### Community 64 - "NetEventType"
Cohesion: 0.07
Nodes (16): snapNetMonitor, NetEvent, NetEventType, NetRecord, BinariesSummary(), TestBinariesSummary(), FormatAddr(), NetEventTypes() (+8 more)

### Community 65 - "unix_connect.c"
Cohesion: 0.40
Nodes (13): attack_connect(), attack_dgram(), attack_uring(), denied(), fill_addr(), handoff(), hold(), main() (+5 more)

### Community 66 - "IntegrationSuite"
Cohesion: 0.08
Nodes (14): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest, IntegrationSuite (+6 more)

### Community 67 - "IntegrationSuite"
Cohesion: 0.21
Nodes (6): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent()

### Community 68 - "GuardEvent"
Cohesion: 0.13
Nodes (8): holdHeadless(), HeadlessLine(), snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestLocalizePaths(), GuardEvent, logBacklogDenial()

### Community 69 - "runNetworkGuard"
Cohesion: 0.11
Nodes (22): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runHeadless(), runNetworkGuard(), runTUI(), WithSuspectsFrom() (+14 more)

### Community 70 - "wipe"
Cohesion: 0.12
Nodes (18): deprovisionKind, isRegularFileTarget(), classifyDeprovision(), classifyDeprovisionErr(), isLockedRegularFileErr(), isNotEncryptedErr(), newBoundedKeyFn(), readKey() (+10 more)

### Community 71 - "LaunchRule"
Cohesion: 0.17
Nodes (14): binaryLaunchRules(), launchRulesByPath(), TestLaunchRulesByPath(), argvEntry, LaunchRule, BuiltinAllowedEnv(), BuiltinLaunchRules(), canonLaunchRules() (+6 more)

### Community 72 - "planUpdaters"
Cohesion: 0.14
Nodes (18): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+10 more)

### Community 73 - "fcntl"
Cohesion: 0.10
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 74 - "filevault_test.go"
Cohesion: 0.17
Nodes (28): looksLikeFileVaultRecord(), sealFileVault(), inode(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestEnsureRecoverySidecarPlaceholder(), TestIsEncryptedDispatchesToFileVaultForRegularFiles(), TestIsFileVaultRecordShape() (+20 more)

### Community 75 - "Vault"
Cohesion: 0.14
Nodes (15): encryptDirectories(), secureResources(), verifyEncryptionState(), confirmRunPrereq(), resolveFilesystemPrereqs(), askEncryption(), TestAskEncryptionSkipsNeedEncryptionFalse(), decryptDirectories() (+7 more)

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "logLine"
Cohesion: 0.14
Nodes (10): logHook, logLine, journalLine(), TestJournalLine(), logRow(), logStyle(), formatEntry(), sortedKeys() (+2 more)

### Community 78 - "NetGuardEvent"
Cohesion: 0.18
Nodes (6): snapNetGuard, NetGuardEvent, protoString(), netGuardEventLine, netGuardEventMsg, netGuardModel

### Community 79 - "process_vm_readv.c"
Cohesion: 0.21
Nodes (11): copy(), is_verb(), main(), read_once(), usage(), dump_memory(), find_all_rw_regions(), find_heap() (+3 more)

### Community 80 - "os.File"
Cohesion: 0.04
Nodes (81): multicallAdmissible(), scriptIdentity(), checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun() (+73 more)

### Community 81 - "controlServer"
Cohesion: 0.15
Nodes (11): readConfigPut(), authSelfPeer(), readClientLines(), startControlServer(), streamEvents(), TestEditControlSessionEndIdempotent(), TestRecordFailureLockout(), configEditor (+3 more)

### Community 83 - "rcUser"
Cohesion: 0.13
Nodes (39): rcTarget, TestHasRCBlocksTrackEnsureAndRemove(), TestRemoveSSHAgentShellEnvKeepsSSHConfig(), applyRC(), bunTargets(), EnsureAddKeysToAgent(), EnsureBunLauncherEnv(), ensureRCBlock() (+31 more)

### Community 84 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.09
Nodes (31): BinaryStat, placedWalk, treeJudge, binaryStatOf(), binaryStatOf(), StatBinary(), StatBinaryFile(), confineBelow() (+23 more)

### Community 85 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 86 - "exe_supersede.h"
Cohesion: 0.22
Nodes (14): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free() (+6 more)

### Community 87 - "loginServer"
Cohesion: 0.06
Nodes (29): markReady(), markReadyWithLogin(), notifySystemdReady(), hasSentinel(), makeRuntimeDir(), openUserRuntimeDir(), plantSentinel(), runtimePathsFor() (+21 more)

### Community 88 - "RunWith"
Cohesion: 0.13
Nodes (18): installGUI(), needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), deployGUI(), stagedBinary (+10 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "expandPlaceholders"
Cohesion: 0.14
Nodes (7): expandPlaceholders(), TestExpandWhitelist(), TestExtraWatchPathsForSkipsMissing(), LibDirGlob, TrustGlob, reservedGlobs(), splitTrustGlob()

### Community 92 - "SanitizeText"
Cohesion: 0.07
Nodes (24): newBinaryVetter(), ownerPaths(), refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), appRunning(), libDirSig(), startMaintainers() (+16 more)

### Community 93 - "netinject.c"
Cohesion: 0.12
Nodes (15): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+7 more)

### Community 94 - "README.md"
Cohesion: 0.32
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 95 - "NetGuard"
Cohesion: 0.06
Nodes (23): lsmHook, foreignPinnedLink(), PumpRingbuf(), reasonText(), eventTypeKey(), mcNameKey(), statInodeKey(), codeVouched() (+15 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "time.Duration"
Cohesion: 0.14
Nodes (14): procStats, statsSample, duration(), cpuTime(), humanBytes(), rssBytes(), New(), recv() (+6 more)

### Community 100 - "RefreshSection"
Cohesion: 0.08
Nodes (33): logRefreshChange(), sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange() (+25 more)

### Community 101 - "btrfsLayoutFrom"
Cohesion: 0.17
Nodes (9): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), btrfsTypes(), TestBtrfsLayoutFromNestedMembers() (+1 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "procstats.go"
Cohesion: 0.30
Nodes (9): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+1 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "ptrace_seize.c"
Cohesion: 0.20
Nodes (11): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file(), main(), peek_buffer(), run_tracer() (+3 more)

### Community 106 - "recordops.c"
Cohesion: 0.36
Nodes (15): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+7 more)

### Community 107 - "runWizard"
Cohesion: 0.10
Nodes (26): askBunTmpdirUsers(), gatherPerUserSetup(), setupBunTmpdir(), addedResources(), applyDiffAdditions(), TestAddedResourcesSkipsAlreadyGuardedSSH(), checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled() (+18 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "BusClient"
Cohesion: 0.19
Nodes (8): BusClient, AttestBusClient(), busClientKey(), execStamp(), imageOf(), pidfdPid(), procIdentity(), suspectLabel()

### Community 110 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.11
Nodes (25): class, pos, span, classOf(), column(), dedent(), Edit(), leadingIndent() (+17 more)

### Community 111 - "hashFileShared"
Cohesion: 0.23
Nodes (8): HashStore, sharedHashEntry, claimHash(), hashFileShared(), hashPathShared(), publishHash(), settled(), unchangedSince()

### Community 112 - "UI"
Cohesion: 0.09
Nodes (20): showChangelog(), channelOf(), confirmUpdate(), isTerminal(), Latest(), RunWith(), TestConfirmUpdateRoutesToActiveUI(), TestConfirmUpdateWithoutTerminalOrUIRefuses() (+12 more)

### Community 113 - "guardUnitTest"
Cohesion: 0.06
Nodes (16): addErrKind, bpfGuardEvent, classifyAddErr(), commMatchesGuardedBinary(), decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel() (+8 more)

### Community 115 - "ConfinedEntry"
Cohesion: 0.24
Nodes (21): ConfinedEntry(), HashStats(), hashedReused(), resetSharedHashes(), TestConcurrentCallersHashOnce(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash(), TestPersistentStoreFilledAfterSyncAndHit() (+13 more)

### Community 116 - "RunCmd"
Cohesion: 0.06
Nodes (51): offerDaemonStart(), runDiffCatalog(), buildBinaryIfNeeded(), certifyInstall(), installBinaryOnly(), markInstallComplete(), installServices(), mustCwd() (+43 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "classify"
Cohesion: 0.19
Nodes (14): call(), TestClassify(), TestClassifyCreateItemReplace(), TestClassifyItemTypeIsSchema(), TestClassifyOpenCallsExact(), arg(), attributesOf(), classify() (+6 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 123 - "newTestTarget"
Cohesion: 0.22
Nodes (17): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+9 more)

### Community 124 - "inet_connect.c"
Cohesion: 0.26
Nodes (15): attack_connect(), attack_fastopen(), attack_uring(), denied(), dial(), fill_addr(), handoff(), listen_on() (+7 more)

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "Config"
Cohesion: 0.04
Nodes (65): ServeConfig, openBinaryVetter(), loadDaemonConfig(), makeReloadHandler(), newDaemonModel(), noteReloadedSettings(), runDaemonUI(), runServedTUI() (+57 more)

### Community 127 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.13
Nodes (15): App maintainers (Steam runtime updates), Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials (+7 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "vaultfs.go"
Cohesion: 0.09
Nodes (26): entryOf(), DefaultNewFileMode(), FileModeToUnixMode(), Vault, IsBinaryContent(), IsOctalMode(), Open(), ParseMode() (+18 more)

### Community 130 - "KernelDev"
Cohesion: 0.25
Nodes (11): fdStat, physicalParent(), statFD(), syncBusBrokers(), TagRelayBus(), KernelDev(), BtrfsDevices(), btrfsSysfsDevices() (+3 more)

### Community 131 - ".admitGated"
Cohesion: 0.18
Nodes (5): placedKey(), Caller, Role, withSchema(), Judge

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "Monitor"
Cohesion: 0.18
Nodes (4): NewMonitor(), putBtrfsLayout(), Monitor, watchedInode

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 144 - "liveSession"
Cohesion: 0.15
Nodes (12): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries(), holdForward() (+4 more)

### Community 145 - "batchInodeRows"
Cohesion: 0.24
Nodes (11): inodeBatchFn, appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches(), TestBatchInodeRows(), walkInodeBatches(), countInodeRows() (+3 more)

### Community 146 - "Ledger"
Cohesion: 0.07
Nodes (21): candidates(), confirmAll(), describe(), record(), run(), applySchema(), Ledger, Pending (+13 more)

### Community 147 - "serveWebSocket"
Cohesion: 0.21
Nodes (8): configureWebSocket(), sameOrigin(), serveWebSocket(), TestSameOrigin(), TestSnapshotHubSingleViewerAndCoalescing(), writeSnapshot(), writeSnapshotLoop(), snapshotHub

### Community 148 - "Daemon configuration"
Cohesion: 0.14
Nodes (13): `[inspectors]`, Daemon configuration, Daemon settings, Generic Electron apps, Launch rules, Library trust, Port sections, Process inspectors (+5 more)

### Community 149 - "runMonitor"
Cohesion: 0.18
Nodes (15): DrainUntilSignal(), BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair() (+7 more)

### Community 150 - ".existingWhitelist"
Cohesion: 0.09
Nodes (23): confSafeMatch(), BinaryRule, homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed() (+15 more)

### Community 152 - "runNetworkMonitor"
Cohesion: 0.10
Nodes (13): ParseEventsFlag(), ParseNetEventsFlag(), parseTypesFlag(), validateForwardEvents(), runHeadless(), runNetworkMonitor(), runTUI(), ProtocolString() (+5 more)

### Community 154 - "find.go"
Cohesion: 0.11
Nodes (27): findButton, findAll(), lowerRunes(), NewFinder(), newFindInput(), overlay(), runesEqual(), findEditor() (+19 more)

### Community 155 - "runGenKey"
Cohesion: 0.20
Nodes (11): confirmMasterKeyOverwrite(), runGenKey(), runOneShotMode(), ensureInstalledBinary(), ensureMasterKey(), prepareInstallation(), TestEnsureInstalledBinary(), checkKeyLen() (+3 more)

### Community 156 - "io/fs.FileInfo"
Cohesion: 0.31
Nodes (10): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), fchownFromInfo(), futimesFromInfo(), readDirNames() (+2 more)

### Community 157 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 158 - "Installation"
Cohesion: 0.17
Nodes (12): Build from source, Bun-based apps (opencode), Docker, Facultatives, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer (+4 more)

### Community 159 - "os.FileInfo"
Cohesion: 0.19
Nodes (5): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata()

### Community 161 - ".Layout"
Cohesion: 0.44
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 162 - "update_test.go"
Cohesion: 0.09
Nodes (19): parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFilterChannel(), testKeyPair(), TestNewerThanInstalled(), TestNewerThanInstalledAcrossTagFormats() (+11 more)

### Community 163 - "IntegrationSuite"
Cohesion: 0.23
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 164 - "absPath"
Cohesion: 0.07
Nodes (8): IntegrationSuite, IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath(), TestMain(), helperChild(), TestMain()

### Community 165 - "feed.go"
Cohesion: 0.05
Nodes (43): actionBar, cell, column, dataTable, eventFeed, feedAction, feedSpec, networkPage (+35 more)

### Community 166 - "shell"
Cohesion: 0.05
Nodes (51): fakeLink, pageDef, pageID, reel, resizeWatch, shell, sizedDialog, snapEngines (+43 more)

### Community 167 - "runDaemon"
Cohesion: 0.07
Nodes (24): drain(), newCatalogRefresher(), readEvents(), swapFile(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestDebouncerPendingCoversWaitAndRun(), newControlManager() (+16 more)

### Community 168 - "tempGrantSetup"
Cohesion: 0.33
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 169 - "changelog.go"
Cohesion: 0.25
Nodes (17): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), compareStableVersions(), filterChannel(), findUpdate(), latestTime() (+9 more)

### Community 170 - "SystemTrusted"
Cohesion: 0.25
Nodes (8): rootOwnedInode(), rootOwnedSafe(), SystemTrusted(), systemTool(), archiveMembers(), decompressed(), tarMembers(), TestArchiveMembersIgnoresNonArchives()

### Community 171 - "hashCache"
Cohesion: 0.19
Nodes (7): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache

### Community 172 - "Limitations"
Cohesion: 0.18
Nodes (11): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Secrets served over IPC (+3 more)

### Community 174 - "parseDaemonEvents"
Cohesion: 0.15
Nodes (6): IntegrationSuite, IntegrationSuite, IntegrationSuite, parseDaemonEvents(), daemonEvent, IntegrationSuite

### Community 176 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 177 - "plan_test.go"
Cohesion: 0.40
Nodes (11): scanLib(), fixture(), TestApplyRefusesAStageChangedAfterValidation(), TestApplyWritesValidatedContent(), TestScanAndValidateAcceptsLibDirContent(), TestScanRefusesSpecialFiles(), TestScanUntouchedLibDir(), TestTouches() (+3 more)

### Community 178 - "sources"
Cohesion: 0.20
Nodes (5): TestGeneratedLibraryBlockRoundTrips(), lstatOf(), outsideHardlinks(), libDirDev(), sources

### Community 180 - "fadescroll_test.go"
Cohesion: 0.32
Nodes (10): newTestFadeScroll(), settle(), TestFadeBarDragEndsGlide(), TestFadeBarDragMovesContent(), TestFadeBarTapPages(), TestFadeScrollBarsHiddenWhenContentFits(), TestFadeScrollBothDirections(), TestFadeScrollClampsOffset() (+2 more)

### Community 181 - "fileedit_symlink_test.go"
Cohesion: 0.42
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 182 - "runView"
Cohesion: 0.07
Nodes (9): daemonPage, logGrid, runHost, runView, stageItem, stageState, newLogGrid(), captureLogs() (+1 more)

### Community 186 - "Binary trust"
Cohesion: 0.22
Nodes (9): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch rules, Launch scan, Multicall binaries, The dynamic loader (+1 more)

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 190 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "github.com/spf13/cobra.Command"
Cohesion: 0.11
Nodes (27): AddHeadlessFlag(), AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+19 more)

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.40
Nodes (5): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog, What an entry can't cover

### Community 193 - "run.go"
Cohesion: 0.36
Nodes (8): planFrom(), newStage(), Run(), runSandbox(), tail(), TestCopyScriptHoldsOnlyTheConfirmedContent(), Spec, stage

### Community 194 - "os.FileMode"
Cohesion: 0.09
Nodes (31): TestDaemonUnitWithMetadataOutput(), dupFD(), installFileAs(), installSSHAgent(), installSSHAgentEnv(), installUnit(), readAtNoFollow(), TestInstallSSHAgentRefusesSymlinkedUnitDir() (+23 more)

### Community 195 - "Existing"
Cohesion: 0.18
Nodes (8): existingSection(), guardedMatch(), overlaps(), refusedPath(), sectionFor(), validateTarget(), Existing, ExistingRule

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "bufio.Reader"
Cohesion: 0.19
Nodes (12): TestReadGrantRequestConfig(), readAuthRequest(), readBareRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest(), TestReadGrantRequestForward() (+4 more)

### Community 202 - "SecretStoreNotes"
Cohesion: 0.39
Nodes (7): dockerKeyringHelper(), gitCredentialHelpers(), readOwnedRegular(), SecretStoreNotes(), TestReadOwnedRegular_RefusesSymlinkAndOtherOwner(), TestSecretStoreNotes(), WarnSecretStores()

### Community 203 - "compileGlobNames"
Cohesion: 0.17
Nodes (9): GlobChild, statFunc, compileGlobNames(), globText(), resolveBits(), resolveChildren(), statKey(), TestCompileGlobNames_RepeatedPatternOrsBits() (+1 more)

### Community 204 - "GuardInodeKey"
Cohesion: 0.06
Nodes (25): portOwner, refusedReplacements, supersededBinary, TestNearestRootClaims(), TestSetInspectorsBounded(), nearestRootClaims(), deleteInoKeys(), Inspector (+17 more)

### Community 206 - "monitor"
Cohesion: 0.07
Nodes (33): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks (+25 more)

### Community 209 - "fetchReleases"
Cohesion: 0.28
Nodes (7): downloadFile(), downloadReleaseFiles(), fetchReleases(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleases(), TestFetchReleasesHTTPError()

### Community 210 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 211 - "app-listener documentation"
Cohesion: 0.33
Nodes (6): app-listener documentation, Command tree, Getting started, Modes, Operating, The daemon in depth

### Community 212 - ".populateMaps"
Cohesion: 0.10
Nodes (7): rootHandle, openRoot(), TestNoteRootMountFollowsMountID(), deleteResKeys(), guardModeKey(), TestReadOnlyWalkToleratesDanglingSymlinks(), releaseRootTrust()

### Community 213 - "WithBottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 214 - "UnifiedDiff"
Cohesion: 0.13
Nodes (19): checkAndApply(), refreshCatalogConfig(), installConfig(), confirmUnit(), ConfirmOverwrite(), ShowDiff(), TestUnifiedDiff(), TestUnifiedDiffIdentical() (+11 more)

### Community 216 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 217 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 218 - "runBackend"
Cohesion: 0.36
Nodes (6): reserveStdout(), runBackend(), LiveModeAvailable(), requireRoot(), runEditProtected(), runLive()

### Community 219 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 220 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 221 - "find"
Cohesion: 0.22
Nodes (16): deletePostBackups(), listBackups(), restoreBackups(), offerBackupCleanup(), TestMovedFlagError(), movedFlagError(), Delete(), Find() (+8 more)

### Community 222 - "ReadInstallCert"
Cohesion: 0.36
Nodes (8): InstallCert, ReadInstallCert(), RequireNoInstallCert(), TestInstallCertRoundTrip(), TestInstallCertSymlinkRefuses(), TestInstallCertUnreadableRefuses(), withCertPath(), OpenRegularNoFollow()

### Community 223 - "listFacultatives"
Cohesion: 0.12
Nodes (26): installFacultatives(), installSSHAgentFacultative(), listFacultatives(), runFacultatives(), facultativeUser(), TestInstallSSHAgentFacultative(), TestListFacultativesIgnoresDaemonConf(), TestListFacultativesMarksGuardedSSH() (+18 more)

### Community 224 - "editView"
Cohesion: 0.12
Nodes (7): codeEntry, editView, Entry, fileSize(), joinRel(), newCodeEntry(), relOf()

### Community 225 - "supersede_probe.c"
Cohesion: 0.76
Nodes (6): answer(), main(), read_ok(), read_secret(), run_child(), spin()

### Community 226 - "sources_test.go"
Cohesion: 0.38
Nodes (4): newSources(), TestSourcesHoldFilesAndArchiveMembers(), TestSourcesRefuseOtherFilesystem(), writeTar()

### Community 228 - "mountTable"
Cohesion: 0.13
Nodes (10): MountInfo, mountinfoUnitTest, MountByID(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo(), userMountable() (+2 more)

### Community 229 - "deploy"
Cohesion: 0.15
Nodes (24): installHooks(), staged(), applyUpdate(), archAssets(), assetsFor(), checkUpdatePreconditions(), deploy(), downloadAndVerify() (+16 more)

### Community 231 - "dialLiveSession"
Cohesion: 0.07
Nodes (29): auditAfterEdit(), auditFindings(), dialLiveSession(), applyEdit(), editAgain(), fetchConfig(), putConfig(), runEditConfig() (+21 more)

### Community 232 - "network monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network monitor: see what a program talks to

### Community 233 - "inode_dev"
Cohesion: 0.43
Nodes (4): ctx_ptr(), inode_dev(), sb_dev(), guard_inode_free()

### Community 234 - "mutedText"
Cohesion: 0.07
Nodes (42): checkGrid, counterStrip, filesPage, homePage, recordForm, targetView, newBinaryList(), newBinaryRow() (+34 more)

### Community 235 - "NewDumpHook"
Cohesion: 0.24
Nodes (11): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+3 more)

### Community 236 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 237 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 238 - "GenerateConf"
Cohesion: 0.05
Nodes (61): TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), discordCatalogEntry(), groupedDiscordConf(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic(), TestCollectFilesystemPrereqsSkipsRegularFiles(), TestGroupedSectionAddressedByEncryptionRoot() (+53 more)

### Community 239 - "run"
Cohesion: 0.27
Nodes (10): fetchRelease(), InstalledGUIVersion(), InstalledVersion(), readInstalledVersion(), run(), runUpdate(), syncGUI(), TestFetchReleaseByTag() (+2 more)

### Community 240 - "shQuote"
Cohesion: 0.17
Nodes (5): IntegrationSuite, IntegrationSuite, overlayPlanter(), shQuote(), IntegrationSuite

### Community 242 - "NewEventFanout"
Cohesion: 0.24
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 243 - "maintainer/plan.go"
Cohesion: 0.31
Nodes (7): isWhiteout(), opaque(), pathByte(), pathRecord(), TestPathRecord(), opKind, walker

### Community 244 - "testLimiter"
Cohesion: 0.43
Nodes (8): gateEvent(), TestGateLogLimiterFlushesOnShutdown(), TestGateLogLimiterFoldsAcrossThreadsAndPids(), TestGateLogLimiterFoldsReadRepeats(), TestGateLogLimiterKeysOnTarget(), TestGateLogLimiterNeverFoldsWhatMatters(), testLimiter(), TestNoLogMetadataBlocks()

### Community 245 - "pinstate.go"
Cohesion: 0.09
Nodes (38): configureDaemonLogging(), lockOneRoot(), lockRootRecovering(), relockStaleVaults(), resolveConfigPath(), runLockdown(), ensurePinStateFilePlaceholder(), ensurePlaceholder() (+30 more)

### Community 246 - "linkView"
Cohesion: 0.31
Nodes (5): linkTarget(), newLinkView(), TestCheckLinksJudgesEachTarget(), entryKind, linkView

### Community 247 - "writeWithin"
Cohesion: 0.25
Nodes (8): atomicWriteAt(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync(), writeWithin()

### Community 248 - "checkBtrfsKeys"
Cohesion: 0.29
Nodes (7): checkBtrfsKeys(), ensureBtrfsLayout(), FirstOnBtrfs(), OnBtrfs(), ResolveBtrfsLayout(), TestResolveBtrfsLayoutHostBTF(), TestSuperblockDevMatchesStatOffBtrfs()

### Community 249 - "parseKernelRelease"
Cohesion: 0.67
Nodes (4): ciliumVersion, parseKernelRelease(), TestParseKernelRelease(), unameKernelVersion()

### Community 250 - "parseResize"
Cohesion: 0.33
Nodes (4): decodeEndsAtEOF(), parseResize(), readResizeLoop(), TestParseResize()

### Community 254 - "VouchedDevs"
Cohesion: 0.29
Nodes (7): hostVouchedDevs(), isFuseType(), parseMountDev(), VouchedDevs(), TestParseMountDevMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed()

### Community 255 - "runRefreshCatalogOnly"
Cohesion: 0.40
Nodes (5): applyLiveRefresh(), runRefreshCatalogOnly(), softenAutomatedRefreshErr(), TestApplyLiveRefreshDeliversReload(), TestSoftenAutomatedRefreshErr()

### Community 257 - "copyXattrs"
Cohesion: 0.50
Nodes (4): TestCopyXattrsPreservesUserAttributes(), copyXattrs(), listXattrNames(), splitXattrNames()

### Community 267 - "watchGroup"
Cohesion: 0.07
Nodes (55): watchGroup, addResource(), applyAllowLib(), applyAnyGroupDirective(), applyBlockLine(), applyConfigLine(), applyDirective(), applyElectronApp() (+47 more)

## Knowledge Gaps
- **33 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener/tools/codecheck`, `open_by_handle_at POC` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 682 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **42 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `go_pkg_strings`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `.netIntegrityContainer`, `IntegrationSuite`, `parseDaemonEvents`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `.TestRecord_GuardParity`, `IntegrationSuite`?**
  _High betweenness centrality (0.032) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `go_pkg_strings` be split into smaller, more focused modules?**
  _Cohesion score 0.09264418096005525 - nodes in this community are weakly interconnected._
- **Why does `Guard` connect `Guard` to `go_pkg_strings`, `time.Duration`, `GuardEvent`, `StatInode`, `time.Time`, `GuardInodeKey`, `buildOneGuard`, `DaemonEvent`, `os.File`, `hashFileShared`, `guardUnitTest`, `.populateMaps`, `Resource`, `loginServer`, `sync.Mutex`, `engine`, `TrustGuard`?**
  _High betweenness centrality (0.028) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.01543068879752297 - nodes in this community are weakly interconnected._
- **Why does `monitor()` connect `monitor` to `time.Duration`, `secret_probe/main.go`, `app-listener documentation`, `Installation`?**
  _High betweenness centrality (0.026) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.09441441441441441 - nodes in this community are weakly interconnected._